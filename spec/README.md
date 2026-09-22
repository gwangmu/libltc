# Lifetime Chart (LTC) Format Specification

## Introduction

The Lifetime Chart (LTC) format is a TOML-compatible text format intended for recording a subject's lifetime history in a timeline.

### Versioning

The format follows the [calendar versioning scheme](https://calver.org/), "YY.MM.N", where "YY" and "MM" are the last two digits of the zero-padded release year and the month in the Gregorian calendar, and "N" is a sequential release number starting from "1" each month. The lowest version is `26.09.1`.

### Principle

The LTC format is intended to record the _whole lifetime_ of a subject, so becoming inaccessible in some way (e.g., if an editing tool is discontinued or a file-hosting service closes) potentially means losing the entire history of the subject. To prevent this destructive scenario, the format is designed based on the following principles:

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

The format consists of _LTC objects_ that specify a dedicated aspect of the chart. LTC objects are divided into three categories: common, main, and chart. Common objects represent chart-independent generic concepts and are embedded in other objects. Main objects represent information that makes up the chart. Chart objects represent a chart itself. The default value of objects can be overridden depending on how they are used in the enclosing object.

Each LTC object type matches a specific TOML object type, but other equivalent TOML object types are also valid. For example, an LTC object corresponding to an inline TOML table _may_ be represented as a fully expanded TOML table, too.

### Common Objects

#### Attribute List

 - TOML object: array of strings
 - Default: empty

An attribute list is an array of strings, each representing a colon-separated (`:`) key-value attribute pair. If multiple colons exist in a string, the leftmost colon separates the key-value pair. If no colon exists, the attribute is given an empty value (`""`).

The same attribute key can appear multiple times in an attribute list. In this case, the duplicated attribute key's values are chained together under the same key.

By default, the `Attrs` field of any TOML-table-type object is considered an attribute list of the enclosing object.

#### Time

 - TOML object: inline table
 - Default: `1` for `Day`, `0` for others

A time object consists of six optional integer-typed fields: `Year`, `Month`, `Day`, `Hour`, `Minute`, and `Second`. If no field is specified, the time is considered _unknown_, and the corresponding field in the enclosing object may be omitted in the LTC file.

The timezone may depend on the subject's physical location associated with the enclosing object, but it's always set to UTC for internal representation. Recognized attributes below:

 - `Incremental`: If not the [chart subject](#Subject)'s `StartDate`, the time is added to the chart subject's `StartDate`.
 - `Approx`: The time is approximate.
 - `Untracked`: At the boundary of a certain period, the time outside this period is untracked.
 - `Emphasized`: The time should be emphasized.

#### Name

 - TOML object: inline table
 - Default: empty for all fields

A name object consists of three optional string-typed fields: `First`, `Middle`, and `Last`. If no field is specified, the name is considered _unknown_, and the corresponding field in the enclosing object may be omitted in the LTC file.

#### ID

 - TOML object: string
 - Default: empty

An ID object is a string that represents the _chart-local_ ID of the enclosing object ("chart-local ID": see [ID Qualification](#ID-Qualification) for more). The LTC file shouldn't contain duplicate chart-local IDs regardless of the enclosing object's type. A chart-local ID should only consist of alphanumeric characters (a-z, A-Z, or 0-9), dashes (`-`), and underscores (`_`).

Across the LTC file load/save boundary, a duplicate, invalid, or empty ID object is assigned a new unique chart-local ID, and its old value is added to the enclosing object's attribute list with a key `OldID` (if it was non-empty).

### Main Objects 

Main objects are divided into two _kinds_: chart-local and imported. Chart-local main objects are those contained in the current LTC file. Imported main objects are those imported from other LTC files via [import objects](#Import). 

#### Header

 - TOML object: top-level table

The header object specifies the LTC format version with the `Version` field. (default: `26.09.1`; the lowest version).

#### Setting

 - TOML object: `Setting` table

The setting object specifies basic information about the LTC file. Fields below:

 - `CalendarSystem`: (string) The calendar system that the [time objects](#Time) will use. (default: `"Gregorian"`)
 - `NoteFormat`: (string) The format of notes for [chart-local events](#Event). (default: `"Markdown"`)
 - `Categories`: (array of strings) All categories that appear at least once in any [non-subchart events](#Event). This field is auto-corrected when loading/saving the LTC file. Categories without events are also preserved.

#### Subject

 - TOML object: `Subject` table

The subject object specifies basic information about the LTC file's subject. The subject is primarily a person, but it can also be a non-person, such as a group of people (e.g., race, country, company, friend group, ...) or a time-sensitive event sequence (e.g., global conflict, curriculum, public gathering, ...). Fields below:

 - `Name`: (name object) The name of the subject. (default: name object default)
 - `StartDate`: (time object) The start date of this subject. For a human subject, this is considered their birthday. (default: time object default)
 - `EndDate`: (time object) The end date of this subject. For a human subject, this is considered their date of death. (default: time object default)
 - `Sex`: (string) The _congenital_ sex of this subject. (default: empty)

Note on the _identified_ sex: Because the identified gender can change over time, it's better to specify it as a [period](#Event) in a separate category (e.g., "Identified Gender") rather than as a field in the subject object that lacks the representation capability of the passage of time.

Note on _undefined_ (not _unknown_) `StartDate`: undefined `StartDate`s may be relevant when the timing information of the events should be specified relative to `StartDate`, but it cannot be pinpointed to a specific time (e.g., a tentative travel plan or an academic curriculum). In such a case, `StartDate` may only specify `Day` as `0`.

#### Event

 - TOML object: table in the `Event` table array

An event object is a fundamental object of the LTC file. It describes a specific event or period during the subject's lifetime. Note that _period_ is not a syntactic concept in the LTC format because the boundary between an event and a period is unclear. Instead, the format does not distinguish them and uses the same event object,

Event objects have three _types_: plain, embedding, and subchart. Plain event objects have no external reference (except in their notes). Embedding event objects embed an event of another LTC file. Subchart event objects embed an entire LTC file. Fields below:

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

`StartDate` should be earlier than or equal to `EndDate`; if not, both dates should be marked as unknown, and their old `StartDate` and `EndDate` should be added to the attribute list with the keys `OldStartDate` and `OldEndDate` across the LTC file load/save boundary, respectively. If either `StartDate` or `EndDate` is unknown, the unknown date is auto-calculated to a month before or after the known one. If both `StartDate` and `EndDate` are unknown, the front-end LTC tool should display such events separately and may not display them on the timeline.

Note on the distinction between subchart event objects and [import objects](#Import): An external LTC file embedded via a subchart event object is still a separate LTC file, so the categories in each chart remain separate. In contrast, an external LTC file imported via an import object is _merged_ into the current LTC file, so the imported event objects are included in the same-name category along with chart-local event objects. Recognized attributes below:

 - `ContinuedFrom:<id>`: This event is continued from another event with a qualified ID `<id>`. See [ID Qualification](#ID-Qualification) for a qualified ID.
 - `AmbiguousPeriod`: This event has an ambiguous period overall.

Note on the distinction between `AmbiguousPeriod` and `Approx` start/end dates: An event may set `Approx` start/end dates if they are independently approximate, and/or set the `AmbiguousPeriod` attribute if the temporal information of the entire event (e.g., duration or approximate start/end dates with wide margins) is largely uncertain.

#### Annex

 - TOML object: table in the `Annex` table array

An annex object represents data attached to the LTC file: photos, text snippets, links, etc. Fields below:

 - `ID`: (ID object) The ID of the annex. (default: ID object default)
 - `Title`: (string) The descriptive title of the annex. (default: empty)
 - `Format`: (string) The data Format of the annex. (default: "txt")
 - `Encoding`: (string) The data encoding of the annex. (default: "none")
 - `Data`: (string) The encoded data of the annex. (default: empty)
 - `Note`: (string) The note of the annex. (default: empty)

The `none` encoding performs no encoding. Since the LTC format is text-based, any binary data should be encoded into a text representation before being included in an LTC file. If `Data` is binary before saving to an LTC file, it should be encoded in a base64 format, add old non-empty `Encoding` to the attribute list as the key `OrgEncoding`, and replace `Encoding` with `base64`. 

By default, `Format` is specific to the front-end LTC tool, except `txt` for text data and `png` for PNG image data. The front-end LTC tool should assume unrecognized `Format`s (including an empty `Format`) as `txt` and report it to users.

#### Import

 - TOML object: table in the `Import` table array

An import object declares an external LTC file to merge into the current LTC file, adding the event objects inside to the same-name category or creating new categories if they don't already exist. See [referencing](#Referencing) for the nested import limit. Fields below:

 - `ID`: (ID object) The ID of the import. (default: ID object default)
 - `Link`: (string) The URI to an import-target LTC file. See [referencing](#Referencing) for a valid URI. (default: empty)
 - `StartDate`: (time object) The start date of the import. (default: time object default)
 - `EndDate`: (time object) The end date of the import. (default: time object default)
 - `Categories`: (array of strings) Categories to import. (default: all)
 - `Note`: (string) The note of the import. (default: empty)

`StartDate` and `EndDate` act as a _period mask_ for imported events. Specifically, after treating unknown `StartDate` and `EndDate` the same way as other [event objects](#Event),

 - For the event objects with unknown `StartDate` and `EndDate`, they are imported unconditionally. 
 - For the event objects that ended before `StartDate` or started after `EndDate`, they are not imported.
 - For the event objects that started after `StartDate` and ended before `EndDate` (inclusive), they are imported as they are.
 - For the event objects that started between `StartDate` and `EndDate` (exclusive) but ended after `EndDate`, their `EndDate` is corrected to the import object's `EndDate` with an `Untracked` attribute.
 - For the event objects that ended between `StartDate` and `EndDate` (exclusive) but started before `StartDate`, their `StartDate` is corrected to the import object's `StartDate` with an `Untracked` attribute.
 
Unknown `StartDate`s or `EndDate`s in the import object are regarded as infinite past or future, respectively. `Categories` also acts as a _category mask_ for imported events, meaning only the event objects in specified categories are imported. Recognized attributes below:

 - `ExcludeCategory:<name>`: don't import the events in the category `<name>`. This category will not be imported even if it is specified in `Categories`.

Note: use [embedding event objects](#Event) to import individual event objects.

### Chart Object

 - TOML object: the entire TOML file

A chart object represents the entire chart described in an LTC file, whose name ends with the extension `.ltc`. An _empty_ chart is defined as a chart with empty main objects.

### ID Qualification

There are two types of IDs: _chart-local_ and _qualified_. The chart-local ID is the one directly specified in the object. The qualified ID is the chart-local ID prefixed with the IDs of the objects through which the corresponding object is included (i.e., imported or embedded) in the current LTC file. The IDs in a prefix are separated by slashes (`/`), and the object IDs with fewer nesting levels should come first. For chart-local (i.e., not included) objects, the prefix is empty. Some examples below:

 - For an event object `e001` that was chart-local, the qualified ID is `e001`.
 - For an event object `e002` that was imported through an import object `i001`, the qualified ID is `i001/e002`.
 - For an event object `e003` that was imported through an import object `i002`, which in turn was imported through an import object `i001`, the qualified ID is `i001/i002/e003`.
 - For an event object `e004` that was embedded through a subchart event object `e999`, the qualified ID is `e999/e004`.

The example above describes only the qualified IDs of event objects, but the same applies to any objects with chart-local IDs (e.g., annex and import objects).

### Referencing

An LTC object with a chart-local ID has both _a local URI_ and _a web URI_. The local URI is either a relative (to the current LTC file's directory) or an absolute filesystem path of the LTC file, plus `.obj/<id>` at the end, where `<id>` is the object's qualified ID. For example, given an LTC file at `/home/john/chart.ltc` that contains an event object `e001`, the local URIs of the chart and the event are `/home/john/chart.ltc` and `/home/john/chart.ltc.obj/e001`, respectively. A Local URI may only specify the qualified ID of the reference object (i.e., may omit the LTC file's filesystem path and `.obj/` at the beginning) if it belongs to the same LTC file. For example, the qualified ID `i001/e002` is the same as the local URI `/home/john/chart.ltc.obj/i001/e002` inside the LTC file `/home/john/chart.ltc`.

The web URI is a web address to an LTC file, combined with an HTML query key `id` at the end. For example, if `https://myltc.com/john` is a web address to an LTC file, the web URI of the event object `e001` in such a file is `https://myltc.com/john?id=e001`. 

In [event notes](#Event), either of the URIs can be used as the _source path/address_ when creating a reference to an object (as a link) or embedding an image. Embedding images directly via a path or a web address is highly discouraged (e.g., `![](/home/john/image.png)` or `![](https://example.com/image.png)`; if any such cases are discovered, the front-end LTC tool should report them and provide an option to include such images as [annex objects](#Annex) across the LTC file load/save boundary.

By default, nested references are limited to 10 times, but the front-end LTC tool may adjust this. If the nested reference exceeds the limit, the front-end LTC tool should report this and treat the final referenced object as an empty object of the same type.
