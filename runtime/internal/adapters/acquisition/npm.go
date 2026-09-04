package acquisition

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/user/orchestra-v3/internal/resources"
	"github.com/user/orchestra-v3/internal/runner"
)

// BannedGlobalFlags defines CLI flags that trigger immediate programmatic rejection
var BannedGlobalFlags = []string{
	"-g",
	"--global",
	"-global",
	"--location=global",
	"-location=global",
	"-g=true",
	"--global=true",
	"-global=true",
}

// DefaultPackageAliases maps canonical resource IDs to exact NPM package specifiers
var DefaultPackageAliases = map[string]string{
	"r3f":            "@react-three/fiber",
	"drei":           "@react-three/drei",
	"vgpu-shaders":   "vgpu",
	"gsap":           "gsap",
	"lenis":          "lenis",
	"motion-dev":     "framer-motion",
	"react-spring":   "@react-spring/web",
	"react-bits":     "react-bits",
	"threeui":        "threeui",
	"trig-js":        "trig-js",
	"vanta":          "vanta",
	"expo-haptics":   "expo-haptics",
	"expo-speech":    "expo-speech",
	"mapcn":          "mapcn",
	"shadergradient": "shadergradient",
	"pretext":        "pretext",
}

// PackageManager represents supported Node.js package managers
type PackageManager string

const (
	PackageManagerPNPM PackageManager = "pnpm"
	PackageManagerYarn PackageManager = "yarn"
	PackageManagerNPM  PackageManager = "npm"
)

// PackageJSON models relevant fields from project package.json
type PackageJSON struct {
	Name             string            `json:"name,omitempty"`
	Version          string            `json:"version,omitempty"`
	Dependencies     map[string]string `json:"dependencies,omitempty"`
	DevDependencies  map[string]string `json:"devDependencies,omitempty"`
	PeerDependencies map[string]string `json:"peerDependencies,omitempty"`
	PackageManager   string            `json:"packageManager,omitempty"`
}

// NPMAdapterOptions configures NPM adapter behavior
type NPMAdapterOptions struct {
	DefaultTimeout time.Duration
	PreferredPM    PackageManager
	SaveDev        bool
	PackageAliases map[string]string
	Locker         ResourceLocker
}

// NPMAdapter implements AcquisitionAdapter for NPM packages
type NPMAdapter struct {
	runner  runner.CommandRunner
	catalog *resources.ResourceCatalog
	options NPMAdapterOptions
	locker  ResourceLocker
}

// NewNPMAdapter creates an NPM adapter instance
func NewNPMAdapter(r runner.CommandRunner, cat *resources.ResourceCatalog, opts *NPMAdapterOptions) *NPMAdapter {
	if r == nil {
		r = runner.NewOSCommandRunner()
	}
	defaultOpts := NPMAdapterOptions{
		DefaultTimeout: 3 * time.Minute,
		PreferredPM:    PackageManagerPNPM,
		PackageAliases: make(map[string]string),
	}
	for k, v := range DefaultPackageAliases {
		defaultOpts.PackageAliases[k] = v
	}
	var locker ResourceLocker
	if opts != nil {
		if opts.DefaultTimeout > 0 {
			defaultOpts.DefaultTimeout = opts.DefaultTimeout
		}
		if opts.PreferredPM != "" {
			defaultOpts.PreferredPM = opts.PreferredPM
		}
		defaultOpts.SaveDev = opts.SaveDev
		if opts.PackageAliases != nil {
			for k, v := range opts.PackageAliases {
				defaultOpts.PackageAliases[k] = v
			}
		}
		if opts.Locker != nil {
			locker = opts.Locker
		}
	}
	if locker == nil {
		locker = NewHybridLocker(filepath.Join(os.TempDir(), "orchestra_npm_locks"))
	}
	return &NPMAdapter{
		runner:  r,
		catalog: cat,
		options: defaultOpts,
		locker:  locker,
	}
}

// SetLocker allows updating or overriding the concurrency locker
func (a *NPMAdapter) SetLocker(l ResourceLocker) {
	if l != nil {
		a.locker = l
	}
}

func (a *NPMAdapter) Name() string {
	return "npm"
}

func (a *NPMAdapter) CanHandle(method string) bool {
	return strings.EqualFold(method, "npm")
}

// CheckGlobalInstallSafety inspects arguments, flags, and destination to enforce anti-global policy
func CheckGlobalInstallSafety(args []string, dest string) error {
	// 1. Check all arguments for banned global flags and patterns
	for i := 0; i < len(args); i++ {
		arg := args[i]
		lower := strings.ToLower(strings.TrimSpace(arg))
		if lower == "" {
			continue
		}

		// A. Check exact banned flags and prefix assignments (e.g. -g, --global, -g=true, -g=false, -g=1, -g=yes)
		for _, banned := range BannedGlobalFlags {
			if lower == banned || strings.HasPrefix(lower, banned+"=") {
				return fmt.Errorf("%w: argument '%s' indicates global installation", ErrGlobalInstallBlocked, arg)
			}
		}

		// B. Check quoted location assignments: --location="global", --location='global'
		if strings.HasPrefix(lower, "--location=") || strings.HasPrefix(lower, "-location=") {
			parts := strings.SplitN(lower, "=", 2)
			val := strings.Trim(parts[1], "\"'")
			if val == "global" {
				return fmt.Errorf("%w: argument '%s' indicates global installation", ErrGlobalInstallBlocked, arg)
			}
		}

		// C. Check two-token location flag: --location global, -location global
		if (lower == "--location" || lower == "-location") && i+1 < len(args) {
			nextVal := strings.Trim(strings.ToLower(strings.TrimSpace(args[i+1])), "\"'")
			if nextVal == "global" {
				return fmt.Errorf("%w: argument '%s %s' indicates global installation", ErrGlobalInstallBlocked, arg, args[i+1])
			}
		}

		// D. Tokenize multi-word strings into discrete words (e.g. "yarn global add")
		words := strings.Fields(lower)
		for _, w := range words {
			cleanWord := strings.Trim(w, "\"'")
			if cleanWord == "global" {
				return fmt.Errorf("%w: 'global' keyword detected in arguments", ErrGlobalInstallBlocked)
			}
			for _, banned := range BannedGlobalFlags {
				if cleanWord == banned || strings.HasPrefix(cleanWord, banned+"=") {
					return fmt.Errorf("%w: argument '%s' indicates global installation", ErrGlobalInstallBlocked, arg)
				}
			}
		}

		// E. Single token with = pointing to system directory: --prefix=/usr, -prefix=/usr
		if strings.HasPrefix(lower, "--prefix=") || strings.HasPrefix(lower, "-prefix=") {
			parts := strings.SplitN(arg, "=", 2)
			targetPath := strings.Trim(parts[1], "\"'")
			if isSystemDirectory(targetPath) {
				return fmt.Errorf("%w: --prefix targets system directory '%s'", ErrGlobalInstallBlocked, targetPath)
			}
		}

		// F. Two tokens in args slice: args[i] is --prefix and args[i+1] is target directory
		if lower == "--prefix" || lower == "-prefix" {
			if i+1 < len(args) {
				targetPath := strings.Trim(args[i+1], "\"'")
				if isSystemDirectory(targetPath) {
					return fmt.Errorf("%w: --prefix targets system directory '%s'", ErrGlobalInstallBlocked, targetPath)
				}
			}
		}

		// G. Multi-word single token containing --prefix (e.g. "--prefix /usr")
		if len(words) >= 2 && (words[0] == "--prefix" || words[0] == "-prefix") {
			targetPath := strings.Trim(words[1], "\"'")
			if isSystemDirectory(targetPath) {
				return fmt.Errorf("%w: --prefix targets system directory '%s'", ErrGlobalInstallBlocked, targetPath)
			}
		}
	}

	// 2. Check destination directory
	if strings.TrimSpace(dest) == "" {
		return fmt.Errorf("%w: empty destination path would default to global or unmanaged context", ErrGlobalInstallBlocked)
	}
	if isSystemDirectory(dest) {
		return fmt.Errorf("%w: destination '%s' targets system-wide directory", ErrGlobalInstallBlocked, dest)
	}

	return nil
}

// isSystemDirectory detects root or system paths
func isSystemDirectory(path string) bool {
	trimmed := strings.TrimSpace(path)
	if trimmed == "" {
		return true
	}
	// Explicit cross-platform slash conversion (filepath.ToSlash is a no-op on Linux for backslashes)
	cleanPath := strings.ReplaceAll(trimmed, "\\", "/")
	norm := strings.ToLower(filepath.ToSlash(filepath.Clean(cleanPath)))
	normNoSlash := strings.TrimRight(norm, "/\\")

	if norm == "/" || norm == "\\" || normNoSlash == "c:" || norm == "c:/" || norm == "c:\\" {
		return true
	}

	systemPaths := []string{
		"/usr",
		"/usr/local",
		"/usr/bin",
		"/usr/local/bin",
		"/etc",
		"/var",
		"c:/windows",
		"c:/windows/system32",
		"c:/program files",
		"c:/program files (x86)",
		"c:/programdata",
	}
	for _, sys := range systemPaths {
		sysClean := strings.ToLower(filepath.ToSlash(filepath.Clean(strings.ReplaceAll(sys, "\\", "/"))))
		sysNoSlash := strings.TrimRight(sysClean, "/\\")
		if norm == sysClean || normNoSlash == sysNoSlash || strings.HasPrefix(norm, sysNoSlash+"/") {
			return true
		}
	}
	// Check user home global npm paths
	if strings.Contains(norm, "/appdata/roaming/npm") || strings.Contains(norm, "/.nvm/") || strings.Contains(norm, "/.npm/") {
		return true
	}
	return false
}

// Acquire executes project-scoped conditional installation of the given resource
func (a *NPMAdapter) Acquire(ctx context.Context, res *resources.Resource, dest string) (*AcquisitionResult, error) {
	start := time.Now()

	// 1. Validate Resource against Canonical Catalog
	if res == nil || res.ID == "" {
		return nil, ErrResourceNotFound
	}
	if a.catalog != nil {
		canonical, found := a.catalog.FindByID(res.ID)
		if !found {
			return nil, fmt.Errorf("%w: resource '%s' not registered in resources.json", ErrResourceNotFound, res.ID)
		}
		if canonical.Status == "REJECTED" {
			return nil, fmt.Errorf("%w: resource '%s' is marked REJECTED by policy", ErrResourceRejected, res.ID)
		}
		res = canonical
	}

	if !a.CanHandle(res.AcquisitionMethod) {
		return nil, fmt.Errorf("%w: resource '%s' requires acquisition method '%s', not npm", ErrResourceNotAllowed, res.ID, res.AcquisitionMethod)
	}

	// 2. Resolve Package Name early
	pkgName := res.ID
	if strings.TrimSpace(res.NpmPackage) != "" {
		pkgName = strings.TrimSpace(res.NpmPackage)
	} else if alias, ok := a.options.PackageAliases[res.ID]; ok {
		pkgName = alias
	}

	// 3. Strict Anti-Global Verification on destination, resource ID, and package name
	if err := CheckGlobalInstallSafety([]string{res.ID, pkgName}, dest); err != nil {
		return nil, err
	}

	// 4. Quarantine Boundary Enforcement
	if err := resources.CheckQuarantineBoundary(dest); err != nil {
		return nil, fmt.Errorf("%w: destination '%s' violates quarantine: %w", ErrQuarantinedDestination, dest, err)
	}

	// 5. Concurrency Control: Serialize workspace installations
	if a.locker != nil {
		lockKey := "workspace:" + filepath.Clean(dest)
		unlock, err := a.locker.Lock(ctx, lockKey)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrResourceLockTimeout, err)
		}
		defer unlock()
	}

	// 5. Inspect Destination Workspace & package.json
	pkgJSONPath := filepath.Join(dest, "package.json")
	pkgData, err := os.ReadFile(pkgJSONPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("%w: %s", ErrPackageJSONNotFound, pkgJSONPath)
		}
		return nil, fmt.Errorf("failed to read package.json: %w", err)
	}

	var pkgJSON PackageJSON
	if err := json.Unmarshal(pkgData, &pkgJSON); err != nil {
		return nil, fmt.Errorf("malformed package.json at %s: %w", pkgJSONPath, err)
	}

	// 6. Conditional Check: Is package already installed?
	alreadyDeclared := false
	declaredVersion := ""
	if v, ok := pkgJSON.Dependencies[pkgName]; ok {
		alreadyDeclared = true
		declaredVersion = v
	} else if v, ok := pkgJSON.DevDependencies[pkgName]; ok {
		alreadyDeclared = true
		declaredVersion = v
	}

	nodeModulesPkgPath := filepath.Join(dest, "node_modules", filepath.FromSlash(pkgName), "package.json")
	nodeModulesExist := false
	installedVersion := declaredVersion

	if stat, err := os.Stat(nodeModulesPkgPath); err == nil && !stat.IsDir() {
		nodeModulesExist = true
		if modData, err := os.ReadFile(nodeModulesPkgPath); err == nil {
			var modJSON struct {
				Version string `json:"version"`
			}
			if json.Unmarshal(modData, &modJSON) == nil && modJSON.Version != "" {
				installedVersion = modJSON.Version
			}
		}
	}

	if alreadyDeclared && nodeModulesExist {
		installedPath := filepath.Dir(nodeModulesPkgPath)
		return &AcquisitionResult{
			ResourceID:        res.ID,
			AdapterName:       a.Name(),
			AcquisitionMethod: "npm",
			PackageName:       pkgName,
			VersionOrSHA:      installedVersion,
			SHA256Hash:        hashFileSHA256(nodeModulesPkgPath),
			ResolvedTarget:    installedPath,
			InstalledPath:     installedPath,
			AlreadyInstalled:  true,
			Duration:          time.Since(start),
			Metadata: map[string]string{
				"status": "already_installed",
				"source": "package.json",
			},
		}, nil
	}

	// 7. Detect Package Manager
	pm, err := a.detectPackageManager(dest, &pkgJSON)
	if err != nil {
		return nil, err
	}

	// 8. Construct Package-Manager Specific Command (STRICTLY PROJECT-SCOPED)
	var cmd runner.Command
	cmd.Dir = dest
	cmd.Timeout = a.options.DefaultTimeout

	switch pm {
	case PackageManagerPNPM:
		cmd.Name = "pnpm"
		if a.options.SaveDev {
			cmd.Args = []string{"add", "-D", pkgName}
		} else {
			cmd.Args = []string{"add", pkgName}
		}
	case PackageManagerYarn:
		cmd.Name = "yarn"
		if a.options.SaveDev {
			cmd.Args = []string{"add", "--dev", pkgName}
		} else {
			cmd.Args = []string{"add", pkgName}
		}
	case PackageManagerNPM:
		fallthrough
	default:
		cmd.Name = "npm"
		if a.options.SaveDev {
			cmd.Args = []string{"install", "--save-dev", pkgName}
		} else {
			cmd.Args = []string{"install", "--save", pkgName}
		}
	}

	// Double-check constructed command against anti-global filter
	if err := CheckGlobalInstallSafety(cmd.Args, dest); err != nil {
		return nil, err
	}

	// 9. Execute via CommandRunner
	runRes, err := a.runner.Run(ctx, cmd)
	if err != nil {
		stdout := ""
		stderr := ""
		if runRes != nil {
			stdout = runRes.Stdout
			stderr = runRes.Stderr
		}
		return nil, fmt.Errorf("%w: %s (stdout: %s, stderr: %s)", ErrInstallationFailed, err, stdout, stderr)
	}

	// 10. Post-Installation Verification
	if stat, err := os.Stat(nodeModulesPkgPath); err == nil && !stat.IsDir() {
		if modData, err := os.ReadFile(nodeModulesPkgPath); err == nil {
			var modJSON struct {
				Version string `json:"version"`
			}
			if json.Unmarshal(modData, &modJSON) == nil && modJSON.Version != "" {
				installedVersion = modJSON.Version
			}
		}
	}

	installedPath := filepath.Join(dest, "node_modules", filepath.FromSlash(pkgName))
	installedPkgJSON := filepath.Join(installedPath, "package.json")
	return &AcquisitionResult{
		ResourceID:        res.ID,
		AdapterName:       a.Name(),
		AcquisitionMethod: "npm",
		PackageName:       pkgName,
		VersionOrSHA:      installedVersion,
		SHA256Hash:        hashFileSHA256(installedPkgJSON),
		ResolvedTarget:    installedPath,
		InstalledPath:     installedPath,
		AlreadyInstalled:  false,
		Duration:          time.Since(start),
		ExecutedCommand:   cmd.Name + " " + strings.Join(cmd.Args, " "),
		Output:            runRes.Stdout,
		Metadata: map[string]string{
			"package_manager": string(pm),
			"exit_code":       fmt.Sprintf("%d", runRes.ExitCode),
		},
	}, nil
}

// detectPackageManager resolves package manager based on package.json, lockfiles, and PATH
func (a *NPMAdapter) detectPackageManager(dest string, pkgJSON *PackageJSON) (PackageManager, error) {
	// 1. Check packageManager field in package.json (e.g. "pnpm@8.15.0")
	if pkgJSON != nil && pkgJSON.PackageManager != "" {
		pmStr := strings.ToLower(pkgJSON.PackageManager)
		if strings.HasPrefix(pmStr, "pnpm") {
			if _, err := a.runner.LookPath("pnpm"); err == nil {
				return PackageManagerPNPM, nil
			}
		} else if strings.HasPrefix(pmStr, "yarn") {
			if _, err := a.runner.LookPath("yarn"); err == nil {
				return PackageManagerYarn, nil
			}
		} else if strings.HasPrefix(pmStr, "npm") {
			if _, err := a.runner.LookPath("npm"); err == nil {
				return PackageManagerNPM, nil
			}
		}
	}

	// 2. Check Lockfiles in Workspace
	if _, err := os.Stat(filepath.Join(dest, "pnpm-lock.yaml")); err == nil {
		if _, err := a.runner.LookPath("pnpm"); err == nil {
			return PackageManagerPNPM, nil
		}
	}
	if _, err := os.Stat(filepath.Join(dest, "yarn.lock")); err == nil {
		if _, err := a.runner.LookPath("yarn"); err == nil {
			return PackageManagerYarn, nil
		}
	}
	if _, err := os.Stat(filepath.Join(dest, "package-lock.json")); err == nil {
		if _, err := a.runner.LookPath("npm"); err == nil {
			return PackageManagerNPM, nil
		}
	}

	// 3. User-configured preference
	if a.options.PreferredPM != "" {
		if _, err := a.runner.LookPath(string(a.options.PreferredPM)); err == nil {
			return a.options.PreferredPM, nil
		}
	}

	// 4. PATH Fallback Order: pnpm -> yarn -> npm
	for _, candidate := range []PackageManager{PackageManagerPNPM, PackageManagerYarn, PackageManagerNPM} {
		if _, err := a.runner.LookPath(string(candidate)); err == nil {
			return candidate, nil
		}
	}

	return "", ErrPackageManagerNotFound
}

func hashFileSHA256(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}
