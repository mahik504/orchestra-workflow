# AirLens — college air-quality GUI

status: shipped
kind: college
created: 2026-08-18
path: `C:\projects\DVP`
github: https://github.com/mahik504/AirLens.git

Canonical = `C:\projects\DVP` + GitHub. Folder on disk is still named DVP. Do not rewrite the app.

## The idea

College **Data Visualization using Python** dashboard. Historical CPCB-style CSV for Indian cities, five screens (Overview, City explorer, Comparison, Trends, Map), optional live Open-Meteo EAQI in the sidebar. Not a hiring flagship.

## Who it is for

The course teacher / viva. Not recruiters.

## Surfaces

desktop (FreeSimpleGUI / PySimpleGUI)

## Stack (this product only)

Python, FreeSimpleGUI (PySimpleGUI fallback), Pandas, NumPy, Matplotlib, Seaborn. Live air: Open-Meteo over stdlib `urllib`. No API key. Skip Stitch.

## 3D

none

## What “sexy” means

Lamp-black / rust desk. Charts fill the maximised window. Map click = mark at the point plus that city’s numbers.

## Visual world

none (desktop GUI, not Stitch)

## Motion (earned only)

none

## Must not look like

Generic slate/blue dashboard. “CPCB night desk” copy. A tiny clipped chart in a 920×400 box.

## Quality (honest, within PySimpleGUI)

Shipped: maximised window, custom theme, five screens, CSV vs live EAQI kept separate, India outline under city dots, click card, teacher `REPORT.md` + zip. That is good enough for the assignment. A later PySimpleGUI pass could split the one-file `app.py`, reduce matplotlib redraw flicker, and tidy native widget padding — not a React rewrite, and not worth touching unless the teacher asks.

## Out of scope this round

Rewriting in a web stack. Mixing EAQI with CPCB AQI. Treating this as internship-cv work.

## Open questions

- [ ] Hand-in submitted?
