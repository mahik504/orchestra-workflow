---
name: orchestra-video
description: A rendered mp4 in the product repo that needs one. HyperFrames is the verified tool. Does not load GSAP. Does not install a video suite into the method repo.
---

# Video — 4.1.0

Use this only when the deliverable is a rendered video file. Source method: HyperFrames (`heygen-com/hyperframes`, checked 2026-11-05). The OpenMontage URL from the earlier list returned 404, so it is not an install target.

## Do

1. Keep the composition in the product repo that needs the file.
2. One scene list, one output size, one caption track if the piece has speech.
3. A still frame is the fallback when the render is heavy.

## Do not

- Load `gsap-core` or the scroll route for an mp4.
- Add HyperFrames to this method repository.
- Claim a render succeeded without the output file.
