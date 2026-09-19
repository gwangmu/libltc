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

The format consists of _LTC objects_ that specify a dedicated aspect of the chart. LTC objects are divided into two categories: common and main. Common objects represent chart-independent generic concepts and are embedded in other objects. Main objects represent information that makes up the chart itself. The default value of common objects can be overridden depending on how they are used in the enclosing object.

Each LTC object type matches a specific TOML object type, but other equivalent TOML object types are also valid. For example, an LTC object corresponding to an inline TOML table _may_ be represented as a fully expanded TOML table, too.

### Common Objects

#### Attribute List

 - TOML object: array of strings.
 - Default: empty

An attribute list is an array of strings, each representing a colon-separated (`:`) key-value attribute pair. If multiple colons exist in a string, the leftmost colon separates the key-value pair. If no colon exists, the attribute is given an empty value (`""`).

The same attribute key can appear multiple times in an attribute list. In this case, the duplicated attribute key's values are chained together to the same key.

#### Time

 - TOML object: inline table.
 - Default: `1` for `Day`, `0` for others.

A time object consists of six optional integer-typed fields: `Year`, `Month`, `Day`, `Hour`, `Minute`, and `Second`. If no field is specified, the time is considered _unknown_. The timezone may depend on the subject's physical location associated with the enclosing object, but it's always set to UTC for internal representation.

If an attribute `Incremental` is set, the time is added to the [chart subject](#Subject)'s `StartDate`. The `Incremental` attribute of the chart subject's `StartDate` is ignored. 

#### Name

 - TOML object: inline table.
 - Default: empty for all fields.

A name object consists of three optional string-typed fields: `First`, `Middle`, and `Last`. If no field is specified, the name is considered _unknown_.

### Main Objects 

#### Header

 - TOML object: top-level table.

The header object specifies the LTC format version with the `Version` field. (default: `26.09.1`; the lowest version).

#### Setting

 - TOML object: `Setting` table.

The setting object specifies basic information about the LTC file. Fields below:

 - `CalendarSystem`: (string) The calendar system that the [time objects](#Time) will use. (default: `"Gregorian"`)
 - `NoteFormat`: (string) The format of `Note`s for [events](#Event). (default: `"Markdown"`)
 - `Categories`: (array of strings) All categories that appear at least once in chart-local events. This field is auto-corrected when loading or saving the LTC file if any category appears in chart-local events but not in the array. Categories without events are also preserved.

#### Entity

 - TOML object: `Entity` table.

The entity object specifies basic information about the LTC file's subject. The subject is primarily a person, but it can also be a non-person, such as a group of people (e.g., race, country, company, friend group, ...) or a time-sensitive event sequence (e.g., global conflict, curriculum, public gathering, ...). Fields below:

 - `Name`: (name object) The name of the entity. (default: name object default)
 - `StartDate`: (time object) The start date of this entity. If the entity is a person, the start date is simply their birthday. (default: earliest `StartDate` in chart-local events, or unknown time if no chart-local events exist)
 - `EndDate`: (time object) The end date of this entity. If the entity is a person, the end date is their day of death. (default: unknown time)
 - `Sex`: (string) The _congenital_ sex of this entity. Because the identified gender can change over time, it's better to specify it as a period in its own category. (default: empty)

#### Event

 - TOML object: `Event` table array.

An event object is a fundamental object of the LTC file. It describes a specific event or period during the entity's lifetime. Notice that a _period_ is not a syntactic concept in the LTC format because the boundary between an event and a period is unclear. Instead, the format represents a period as an event object.

Event objects have two key dimensions. On one dimension (_kind_), event objects are classified into _chart-local_ and _imported_. Chart-local event objects are those contained in the current LTC file. Imported event objects are those imported from other LTC files via (import objects)[#Import]. 

On the other dimension (_type_), event objects are classified into _plain_, _embedding_, and _subchart_. Plain event objects directly describe the event in the LTC file. Embedding event objects embed an event from another LTC file. Subchart event objects embed the whole external LTC file.

There is an important distinction between subchart event objects and (import objects)[#Import]. An external LTC file embedded via a subchart event object is still a separate LTC file; the event object may link to the embedded chart, but the categories in each chart remain separate. In contrast, an external LTC file imported via an import object is _merged_ into the current LTC file. As a result, the categories with the same name display both the chart-local and imported event objects.

#### Annex

#### Import
