---
name: ow-video
description: A rendered mp4 in the product repo that needs one. HyperFrames is the verified tool. Does not load GSAP. Does not install a video suite into the method repo.
---

# Video

Use this only when the deliverable is a rendered video file. Source method: HyperFrames (`heygen-com/hyperframes`, checked 2026-11-05). The OpenMontage URL from the earlier list returned 404, so it is not an install target.

## Do

1. Keep the composition in the product repo that needs the file.
2. One scene list, one output size, one caption track if the piece has speech.
3. A still frame is the fallback when the render is heavy.

## Do not

- Load `gsap-core` or the scroll route for an mp4.
- Add HyperFrames to this method repository.
- Claim a render succeeded without the output file.

## Preconditions

- The deliverable is a rendered video file in a product repo.

## Procedure

1. The task's paths and acceptance line win.
2. Write one scene list, one output size, and a caption track if there is speech.
3. Render in the product repo. Keep a still as the fallback.

## Output

The mp4 path, or the still and the reason the render did not finish.

## Verify

The output file exists on disk. A missing file is not a successful render.

## Skip

Skip for page scroll, which is the GSAP route, and for installing a video suite into this method repo.
