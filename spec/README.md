# Lifetime Chart (LTC) Format Specification

## Introduction

The Lifetime Chart (LTC) format is a TOML-compatible text format intended for recording a subject's lifetime history in a timeline.

### Versioning

The format follows the [calendar versioning scheme](https://calver.org/), "YY.MM.N", where "YY" and "MM" are the last two digits of the zero-padded release year and the month in the Gregorian calendar, and "N" is a sequential release number starting from "1" each month. The lowest version is `26.09.1`.

### Principle

The LTC format is intended to record the _whole life_ of a subject, so becoming inaccessible in some way (e.g., if an editing tool is discontinued or a file-hosting service closes) potentially means losing the entire history of the subject, presumably unintentionally. To prevent this destructive scenario, the format is designed based on the following principles:

 - The format should be _human-understandable_ via a text editor. Even without detailed knowledge of the specification, users should be able to deduce the meaning of the contents using a text editor.
 - The format (and the file written in it) should be _always accessible_, otherwise some LTC files may end up being locked into some private business.
 - The format should be _lenient_ so that some deviations from the specification shouldn't invalidate the whole file.
 - The format should be _backward compatible_, otherwise an LTC file may not be accessible if a tool does not support an old format.

The LTC format implements the above principles as follows:

 - The format is based on the [TOML](https://toml.io/en/) format.
 - The format is _open_ (conditioned MIT License), and the license includes an "always-downloadable" condition for files written in it.
 - The format attempts to _auto-correct_ the deviations from the specification in a sensible way, and reports them when it does.
 - A recent version of the format is a _strict superset_ of any preceding versions.

## Specification

The ["sample" directory](./sample) collects the sample LTC file for each version. The sample file is for demonstration purposes only; it is syntactically correct but may be semantically invalid (e.g., mutually incompatible attributes). The _default_ is applied when the corresponding field is unspecified or invalid, and all unrecognized fields and attributes are ignored but preserved across the LTC file load/save boundary. If not specified separately, the specification is based on the lowest version (`26.09.1`).

### LTC Object Taxonomy

The format consists of _LTC objects_ that specify a dedicated aspect of the chart. LTC objects are divided into three categories: common, main, and chart. Common objects represent chart-independent generic concepts and are embedded in other objects. Main objects represent information that makes up the chart itself. Chart objects represent a stand-alone chart. The default value of common objects can be overridden depending on how they are used in the enclosing object.

Each LTC object type matches a specific TOML object type, but other equivalent TOML object types are also valid. For example, an LTC object corresponding to an inline TOML table _may_ be represented as a fully expanded TOML table, too.

### Common Objects

#### Attribute List

 - TOML object: array of strings.
 - Default: empty

An attribute list is an array of strings, each representing a colon-separated (`:`) key-value attribute pair. If multiple colons exist in a string, the leftmost colon separates the key-value pair. If no colon exists, the attribute is given an empty value (`""`).

The same attribute key can appear multiple times in an attribute list. In this case, the duplicated attribute key's values are chained together under the same key.

By default, the `Attrs` field of any TOML-table-type object is considered an attribute list.

#### Time

 - TOML object: inline table.
 - Default: `1` for `Day`, `0` for others.

A time object consists of six optional integer-typed fields: `Year`, `Month`, `Day`, `Hour`, `Minute`, and `Second`. If no field is specified, the time is considered _unknown_, and the corresponding field in the enclosing object may be omitted in the LTC file.

The timezone may depend on the subject's physical location associated with the enclosing object, but it's always set to UTC for internal representation. Recognized attributes below:

 - `Incremental`: If not the [chart subject](#Subject)'s `StartDate`, the time is added to the chart subject's `StartDate`.
 - `Approx`: The time is approximate.
 - `Untracked`: At the boundary of a certain period, the time outside this period is untracked.
 - `Emphasized`: The time should be emphasized.

#### Name

 - TOML object: inline table.
 - Default: empty for all fields.

A name object consists of three optional string-typed fields: `First`, `Middle`, and `Last`. If no field is specified, the name is considered _unknown_, and the corresponding field in the enclosing object may be omitted in the LTC file.

#### ID

 - TOML object: string.
 - Default: empty.

An ID object is a string that represents the _chart-local_ ID of the enclosing object ("chart-local ID": see [ID Qualification](#ID-Qualification) for more). The LTC file shouldn't contain duplicate ID object values regardless of the enclosing object's type. Across the LTC file load/save boundary, a duplicate or empty ID object is assigned a new unique chart-local ID, and its old value is added to the enclosing object's attribute list with a key `OldID` (if it was non-empty).

### Main Objects 

#### Header

 - TOML object: top-level table.

The header object specifies the LTC format version with the `Version` field. (default: `26.09.1`; the lowest version).

#### Setting

 - TOML object: `Setting` table.

The setting object specifies basic information about the LTC file. Fields below:

 - `CalendarSystem`: (string) The calendar system that the [time objects](#Time) will use. (default: `"Gregorian"`)
 - `NoteFormat`: (string) The format of notes for [chart-local events](#Event). (default: `"Markdown"`)
 - `Categories`: (array of strings) All categories that appear at least once in chart-local events. This field is auto-corrected when loading/saving the LTC file. Categories without events are also preserved.

#### Subject

 - TOML object: `Subject` table.

The subject object specifies basic information about the LTC file's subject. The subject is primarily a person, but it can also be a non-person, such as a group of people (e.g., race, country, company, friend group, ...) or a time-sensitive event sequence (e.g., global conflict, curriculum, public gathering, ...). Fields below:

 - `Name`: (name object) The name of the subject. (default: name object default)
 - `StartDate`: (time object) The start date of this subject. For a human subject, this is considered their birthday. (default: time object default)
 - `EndDate`: (time object) The end date of this subject. For a human subject, this is considered their date of death. (default: time object default)
 - `Sex`: (string) The _congenital_ sex of this subject. (default: empty)

Note on the _identified_ sex: Because the identified gender can change over time, it's better to specify it as a [period](#Event) in a separate category (e.g., "Identified Gender") rather than as a field in the subject object that lacks the representation capability of the passage of time.

#### Event

 - TOML object: a table in an `Event` table array.

An event object is a fundamental object of the LTC file. It describes a specific event or period during the subject's lifetime. Note that _period_ is not a syntactic concept in the LTC format because the boundary between an event and a period is unclear. Instead, the format does not distinguish them and uses the same event object,

Event objects have two key dimensions: kind and type. On the _kind_ dimension, event objects are classified into _chart-local_ and _imported_. Chart-local event objects are those contained in the current LTC file. Imported event objects are those imported from other LTC files via [import objects](#Import). 

On the _type_ dimension, event objects are classified into _plain_, _embedding_, and _subchart_. Plain event objects have no external reference (except in their notes). Embedding event objects embed an event of another LTC file. Subchart event objects embed an entire LTC file. Fields below:

 - `ID`: (ID object) The ID of the event. (default: ID object default)
 - `Title`: (string) The descriptive summary ("title") of the event. (default: empty)
 - `Category`: (string) The category of the event. (default: empty)
 - `StartDate`: (time object) The start date of the event. (default: time object default)
 - `EndDate`: (time object) The end date of the event. (default: time object default)
 - `Subchart`: (string) The URI to an embed-target LTC file as a "subchart". See [referencing](#Referencing) for a valid URI. (default: empty)
 - `Embed`: (string) The URI to an embed-target event. See [referencing](#Referencing) for a valid URI. (default: empty)
 - `Note`: (string) The note of the event. (default: empty)

An event object is _embedding-typed_ with a non-empty `Embed` field, _subchart-typed_ with a non-empty `Subchart` field, or _plain-typed_ otherwise. The `Embed` and `Subchart` fields are mutually exclusive; if they both exist, the front-end LTC tool arbitrarily takes one of them and reports that the other was ignored. `Subchart`s can reference the current LTC file, and `Embed`s can reference an event in the current LTC file. See [referencing](#Referencing) for nested references.

For embedding event objects, `StartDate`, `EndDate`, and `Title` are overridden by the embedded event's `StartDate`, `EndDate`, and `Title`, respectively, unless they are unknown in the embedded event. For subchart event objects, `StartDate`, `EndDate`, and `Title` are overridden by the subchart subject's `StartDate`, `EndDate`, and the stringified subchart subject's `Name`, respectively, unless they are unknown in the embedded subchart. `Note` is valid for all event object types. 

Note on the distinction between subchart event objects and [import objects](#Import): An external LTC file embedded via a subchart event object is still a separate LTC file, so the categories in each chart remain separate. In contrast, an external LTC file imported via an import object is _merged_ into the current LTC file, so the imported event objects are included in the same-name category along with chart-local event objects. Recognized attributes below:

 - `ContinuedFrom:<id>`: This event is continued from another event with a qualified ID `<id>`. See [ID Qualification](#ID-Qualification) for a qualified ID.
 - `AmbiguousPeriod`: This event has an ambiguous period overall.

Note on the distinction between `AmbiguousPeriod` and `Approx` start/end dates: An event may set `Approx` start/end dates if they are independently approximate, or set the `AmbiguousPeriod` attribute if the temporal information of the entire event (e.g., duration or approximate start/end dates with wide margins) is largely uncertain.

#### Annex

#### Import

### Chart Object

 - TOML object: the entire TOML file.

TODO: define "empty"

### ID Qualification

### Referencing

By default, nested references are limited to 10 times, but the front-end LTC tool may adjust this. If the nested reference exceeds the limit, the front-end LTC tool should report this and treat the final referenced object as an empty object of the same type.
