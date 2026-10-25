# Validation for 3.3.5

The pin moved. The published tree did not pass its own checks.

| Check | Result on `7543293` |
| --- | --- |
| GitHub CI | Red. [run 40074533564](https://github.com/mahik504/orchestra-workflow/actions/runs/40074533564) died at `go fmt` |
| GitHub Hygiene | Red. [run 40074533531](https://github.com/mahik504/orchestra-workflow/actions/runs/40074533531) crashed with `EISDIR` on the `skills/humanizer` gitlink |
| 3H | Not run. `claim` stays none |

Do not treat the 3.3.5 pin as a green release. 3.3.6 is the patch that fixes these two checks.
