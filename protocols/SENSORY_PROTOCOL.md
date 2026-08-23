# Sensory protocol (optional)

Capabilities: `haptics`, `ui_audio`, `voice_presence`. **Off** unless the product is native/mobile or an assistant (operator HUD). Never decorate CRUD web with whooshes.

If quality is not available, **skip**. Do not ship a robot default.

## Haptics (prefer OS)

| Platform | Resource | Class |
| --- | --- | --- |
| iOS | `UIFeedbackGenerator` / HIG Playing Haptics | SPECIALIST |
| Android | `VibrationEffect` / `HapticFeedbackConstants` | SPECIALIST |
| Expo | `expo-haptics` (wraps the above) | SPECIALIST |

Use for: confirm, reject, PTT start/end, destructive. Not for every tap.

Docs:

- https://developer.apple.com/design/human-interface-guidelines/playing-haptics
- https://developer.android.com/reference/android/os/VibrationEffect
- https://docs.expo.dev/versions/latest/sdk/haptics/

## UI audio

Prefer **system** sounds (`AudioServicesPlaySystemSound` / Android `SoundEffectConstants`) or one authored asset. `expo-av` only if Plan names a file we own.

No random UI pack whooshes. Honor silent switch / ringer. Always a mute path.

## Voice / TTS (operator-HUD)

Order:

1. Platform TTS for accessibility (`AVSpeechSynthesizer`, Android `TextToSpeech`, Expo `expo-speech`) — **accessibility**, not the product voice.
2. Product voice only if Plan names a **high-quality** provider and keys live in **app env**, never the vault.
3. OPTIONAL paid: ElevenLabs / Cartesia **when Plan names them**. Not CORE.

Conversational UX: barge-in, PTT vs always-on already in operator HUD idea.md. Don’t invent Jarvis clones.

**REJECTED:** unnamed “free TTS” voices, autoplay speech on marketing sites.
