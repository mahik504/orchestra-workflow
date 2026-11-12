---
name: ow-ml-experiment
description: Run or report a model experiment with an honest eval. Use for training, fine-tuning, or a dataset gate. Do not invent accuracy.
---

# ML experiment

| Rule | Detail |
| --- | --- |
| Data | Name the dataset and the split before training |
| Eval | Report the metric the split actually produced |
| Repeat | Record seed, command, and commit |
| Soup | Use the fine-tune tool only when the user names it |

If the number is not in a log, say it was not measured. The record names the dataset, the split, the metric, the seed, and the commit. A leaderboard number you did not reproduce stays out of the README.


## Preconditions

- The task names the files or the acceptance line.

## Procedure

1. The task's paths and acceptance line win.
2. Report the command and the number the run printed.

## Output

The change, the command that checked it, and what that command printed.

## Verify

Say the work is done only after that command has been run.

## Skip

Skip when the human says skip orchestra. Do not add a second conductor.
