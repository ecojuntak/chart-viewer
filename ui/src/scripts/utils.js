async function flattenJSONObject(object, parent, options) {
  if (typeof object !== 'object' || Array.isArray(object) || object === null) {
      const key = parent
      const value = object
      const type = Array.isArray(object) ? 'array' : typeof object

      options.push({key: key, value: value, type: type})
      return options
  }

  if (typeof object === 'object' && Object.keys(object).length === 0) {
      const key = parent
      const value = '{}'
      const type = Array.isArray(object) ? 'array' : typeof object

      options.push({key: key, value: value, type: type})
      return options
  }

  const properties = Object.keys(object)

  for (let i = 0; i < properties.length; i++) {
      const name = properties[i]

     await flattenJSONObject(object[name], parent + '.' + name, options)
  }

  return options
}

async function simplifyValues(values) {
  let options = []

  for (let i = 0; i < values.length; i++) {
    const value = values[i]

    const path = value.key.substring(1);

    options.push(path + '=' + value.value)
  }

  return options
}

module.exports = {
  flattenJSONObject,
  simplifyValues
}