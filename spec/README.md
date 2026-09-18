# LTC Specification

## Introduction

The Lifetime Chart (LTC) format is a TOML-compatible text format intended for recording an individual's lifetime history in a timeline.

### Versioning

The format follows the [calendar versioning scheme](https://calver.org/), "YY.MM.N", where "YY" and "MM" are the last two digits of the zero-padded release year and the month in the Gregorian calendar, and "N" is a sequential release number starting from "1" each month. The lowest version is "26.09.1".

### Principle

The LTC format is intended to record the _whole life_ of an individual, so becoming inaccessible in some way (e.g., if an editing tool is discontinued or a file-hosting service closes) potentially means losing the entire history of an individual, presumably unintentionally. To prevent this destructive scenario, the format is designed based on the following principles:

 - The format itself should be _human-understandable_ via a text editor. There will be a format specification, but even without detailed knowledge of the exact specification, users should be able to (easily) deduce the meaning of the contents and possibly reconstruct them manually in an emergency.
 - The format is _as permissive as possible_ in that some deviations from the specification shouldn't invalidate the whole format.

The LTC format implements the above principles as follows:

 - The format is _open_ (MIT License) and _human-readable_ ([TOML](https://toml.io/en/)-compatible).
 - The format attempts to _auto-correct_ the deviations from the specification in a _reasonable_ (not wrong) way, and reports them when it does.

## 
