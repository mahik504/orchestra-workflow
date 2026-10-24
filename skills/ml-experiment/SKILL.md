---
name: ml-experiment
description: Run or report a model experiment with an honest eval. Use for training, fine-tuning, or a dataset gate. Do not invent accuracy.
---

# ML experiment

| Rule | Detail |
| --- | --- |
| Data | Name the dataset and the split before training |
| Eval | Report the metric the split actually produced |
| Repeat | Record seed, command, and commit |
| Soup | Use the fine-tune tool only when the user names it |

If the number is not in a log, say it was not measured.
