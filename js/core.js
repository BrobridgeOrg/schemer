function scanStruct(obj) {
	for (let key in obj) {
		let val = obj[key]
		if (val === undefined) {
			delete obj[key]
		} else if (val == null) {
			continue
		} else if (val instanceof Array) {
			scanStruct(val)
		} else if (typeof val === 'object') {
			scanStruct(val)
		}
	}
}

function main(source) {
	let v = execute(source, script)
	if (!v)
		return null

  // If the result is a proxy, then get the target
  if (v.$target) {
    v = v.$target;
  }

	scanStruct(v)

	return v
}

function normalize(source) {

  let result = {};

  for (let sourceKey in source) {

    let keyParts = sourceKey.split('.');
    let value = source[sourceKey];

    let level = result;
    let count = 0;
    for (let part of keyParts) {

      if (count === keyParts.length - 1) {
        level[part] = value
        break;
      }

      if (!level[part]) {
        level[part] = {}
      }

      level = level[part];
      count++;
    }
  }

  return result;
}

function getUpdates(source, from) {

  let updates = {};

  Object.keys(source).forEach((sourceKey) => {

    // If the source key is the same as the mapping key, then copy the value
    if (sourceKey == from || sourceKey.startsWith(from + '.')) {
      updates[sourceKey] = source[sourceKey];
    }
  });

  return updates;
}

function getValue(source, key) {

  let keyParts = key.split('.');
  let cur = source.$ref;
  for (let part of keyParts) {

    if (!cur[part]) {
      return undefined;
    }

    cur = cur[part];
  }

  return cur;
}

function watch(source, key, cb) {

  let value = getValue(source, key);
  if (value === undefined) {
    return;
  }

  cb(value);
}

function mapping(source = {}, mapping = {}) {

  const target = {};

  // Iterate over the mapping
  Object.entries(mapping).forEach(([key, from]) => {

    let updates = getUpdates(source, from);

    Object.entries(updates).forEach(([sourceKey, value]) => {
      let newKey = key + sourceKey.substring(from.length);
      target[newKey] = value;
    });
  });

  // removed fields
  if (source['$removedFields']) {
    let removedFields = [];
    source.$removedFields.forEach((removedField) => {

      Object.entries(mapping).forEach(([key, from]) => {

        if (removedField == from) {
          removedFields.push(key);
        }
      });
    });

    if (removedFields.length > 0) {
      target.$removedFields = removedFields;
    }
  }

  return target;
}

const proxyHandler = {
  get: function(target, key) {

    switch(key) {
    case '$target':
      return target;
    case '$ref':
      return target.$ref;
    case '$mapping':
      return function(mappingTable) {
        return mapping(target, mappingTable);
      };
    case '$watch':
      return function(key, cb) {
        watch(target, key, cb);
      };
    case '$getValue':
      return function(key) {
        return getValue(target, key);
      };
    case '$getUpdates':
      return function(from) {
        return getUpdates(target, from);
      };
    }

    return target[key];
  }
};


function execute(source, script) {

  let data = Object.assign({}, source, {
    $ref: normalize(source)
  });

  let proxy = new Proxy(data, proxyHandler);

  // Hide properties
  Object.defineProperty(proxy, '$ref', {
    enumerable: false,
    writable: false,
    configurable: true
  });

  let result = script(proxy);
  if (!result) {
    return null;
  }

  return result;
}

module.exports = {
  getUpdates: getUpdates,
  getValue: getValue,
  watch: watch,
  execute: execute,
  mapping: mapping,
};
