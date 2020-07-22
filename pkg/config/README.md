# Config package

The purpose of this package is to hold the configuration related stuff.

## ConfigItem

It holds the information for one configuration item.

- **label** - It is the label of the item. It is displayed on the form screen as the label of the input.
- **description** - It is the description of the item. It is displayed on the form screen in the detailContentBox.
- **valueType** - This enum describes the type of the configuration item. Currently the following types are supported: `int`, `int64`, `float`, `text`, `vector`, `bool`
- **defaultValue** - This is the initial value of the item.
- **currentValue** - This is the current value of the item. If the input value is modified, this value is also modified.
- **key** - This is a uniq string identifier of the form item. This value is used for the maps.

The package provides the `NewConfigItem` function for creating config items. The `SetCurrentValue` updates the current value if the given new value, and the type of the default value are the same.

## Config

This structure holds ConfigItems. It maps the config item keys to config items. The `New` function is responsible for creating a new (empty) config map. The `AddConfig` function can be used to insert a new ConfigItem. Under the hood, it creates the new ConfigItem with the NewConfigItem function. With the `SetCurrentValue` function we can set the current value of the ConfigItem that is mapped to the given key.
