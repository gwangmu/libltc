# LTC Specification

## Introduction

The Lifetime Chart (LTC) format is a TOML-compatible text format intended for recording an individual's lifetime history in a timeline.

### Versioning

The format follows the [calendar versioning scheme](https://calver.org/), "YY.MM.N", where "YY" and "MM" are the last two digits of the zero-padded release year and the month in the Gregorian calendar, and "N" is a sequential release number starting from "1" each month. The lowest version is `26.09.1`.

### Principle

The LTC format is intended to record the _whole life_ of an individual, so becoming inaccessible in some way (e.g., if an editing tool is discontinued or a file-hosting service closes) potentially means losing the entire history of an individual, presumably unintentionally. To prevent this destructive scenario, the format is designed based on the following principles:

 - The format should be _human-understandable_ via a text editor. A format specification will exist, but even without detailed knowledge of it, users should be able to deduce the meaning of the contents using a text editor and, if needed, reconstruct them manually in an emergency.
 - The format should be _as permissive as possible_ in that some deviations from the specification shouldn't invalidate the whole format.
 - The format should be _backward compatible_, otherwise an LTC file may not be accessible if a tool does not support an old format.

The LTC format implements the above principles as follows:

 - The format is _open_ (MIT License with conditions) and _human-readable_ ([TOML](https://toml.io/en/)-compatible).
 - The format attempts to _auto-correct_ the deviations from the specification in a _reasonable_ (not wrong) way, and reports them when it does.
 - A recent version of the format is a _strict superset_ of any preceding versions.

## Specification

The ["sample" directory](./sample) collects the sample LTC file for each version. The sample file is for demonstration purposes only; it is syntactically correct but may be semantically invalid (e.g., mutually incompatible `Attrs`). The _default_ is applied when the corresponding field is unspecified or invalid. If not specified separately, the specification is based on the lowest version (`26.09.1`).

### Common

#### Attributes

Each TOML table may specify specialized attributes via the `Attrs` field. `Attrs` is an array of string-type attributes, each of which can either be a single value or a colon-separated (`:`) key-value pair (if a colon character exists). The key-value pair is separated on the leftmost colon in the string. Unrecognized attributes (i.e., unrecognized attribute values, keys, or values) are order-and-value-preserved across the LTC file load/save boundary. (default: empty)

#### Time

In the LTC file, a time is specified in a predefined TOML table that contains six optional fields: `Year`, `Month`, `Day`, `Hour`, `Minute`, and `Second`. All fields are integer-typed.  The default value is `1` for `Day` and `0` for all others, but it can be overridden depending on the specific usage of the time.

### Tables

#### Top-level

The top-level TOML table specifies the LTC format version with the `Version` field. (default: `26.09.1`; the lowest version).

#### Setting

The `Setting` TOML table specifies basic information about the LTC file.

 - `DisplayLanguage`: (string) The display language of the LTC tool; unrelated to the rest of the file. (default: `"English"`)
 - `CalendarSystem`: (string) The calendar system that the times in this file will use. (default: `"Gregorian"`)
 - `NoteFormat`: (string) The format of `Note`s for [events](#Event). (default: `"Markdown"`)
 - `Categories`: (array of strings) All categories that appear at least once in chart-local events. This field is auto-corrected when loading or saving the LTC file if any category appears in chart-local events but not in the array. Categories without any events will also be preserved.

#### Entity

#### Event

Two kinds of events: chart-local and imported.

Three types of events: direct, embedded, and subchart.

#### Annex

#### Import
