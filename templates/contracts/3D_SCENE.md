# 3D scene

## Purpose

What the scene explains. A scene whose only job is to look expensive does not ship.

## Camera

Position, move, and what it is looking at.

## Lighting

The lights and why each one is there.

## Materials

The surfaces and where the textures come from.

## Loading

What the visitor sees before the scene is ready. A local still, not a spinner.

## Disposal

Geometries, materials, textures, and the renderer are disposed on unmount.

## DPR

The device pixel ratio cap. Default range is 1 through 2.

## Mobile fallback

The still or simpler scene used on a low-power device.

## Reduced motion

The render loop stays stopped.

## WebGL failure

What happens when the context cannot be created or is lost.
