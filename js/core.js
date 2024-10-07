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
	let v = transform(source, script)
	if (!v)
		return null
	scanStruct(v)
	return v
}

function createMappingObject(path) {
  const handler = {
    get(target, prop) {
      if (prop === '$path') {
        return path;
      }

      if (prop === Symbol.toPrimitive) {
        return (hint) => {
          if (hint === 'number') return NaN;
          if (hint === 'string') return '[object Object]';
          return true;
        };
      }

      if (!(prop in target)) {
        target[prop] = createMappingObject();
      }

      return target[prop];
    }
  };

  const proxy = new Proxy({}, handler);

  Object.defineProperty(proxy, '$path', {
    enumerable: false,
    writable: false,
    configurable: true
  });

  return proxy;
}

function trimInvalidKeys(source) {
  const keys = Object.keys(source);
  const filteredKeys = [];

  keys.sort((a, b) => b.length - a.length);

  keys.forEach((key) => {
    if (!filteredKeys.some(k => k.startsWith(key + '.'))) {
      filteredKeys.push(key);
    }
  });

  const result = {};
  filteredKeys.forEach(key => {
    result[key] = source[key];
  });

  return result;
}

function getUpdates(actions) {

  let result = {};

  for (let action in actions) {

    let keys = action.split('.');
    let value = actions[action];

    let level = result;
    let count = 0;
    for (let key of keys) {

      if (count === keys.length - 1) {
        level[key] = value;
        break;
      }

      if (!level[key]) {
        level[key] = {};
      }

      level = level[key];
      count++;
    }
  }

  return result;
}

function createMapping(actions) {

  let result = createMappingObject();

  for (let action in actions) {

    if (!action.includes('.')) {
      continue;
    }

    let keys = action.split('.');
    let value = actions[action];

    let level = result;
    let count = 0;
    for (let key of keys) {

      if (count === keys.length - 1) {
        level[key] = createMappingObject(action);
        break;
      }

      if (!level[key]) {
        level[key] = createMappingObject();
      }

      level = level[key];
      count++;
    }
  }

  return result;
}

function findTargetKey(originKey, mapping) {

  for (let key in mapping) {
    let value = mapping[key];

    if (value.$path === originKey) {
      return key;
    }

    // Next level
    if (typeof value === 'object') {
      let targetKey = findTargetKey(originKey, value);
      if (targetKey) {
        return key + '.' + targetKey;
      }
    }
  }

  return null;
}

function getValueFromStruct(key, struct) {
  let keys = key.split('.');
  let level = struct;
  for (let key of keys) {
    level = level[key];
    if (!level) {
      break;
    }
  }

  return level;
}

function cleanMapping(mapping, updates) {

  for (let key in mapping) {
    let value = mapping[key];

    if (value === undefined) {
      delete mapping[key]
      continue;
    }

    if (typeof value === 'object' && value.$path === undefined) {
      mapping[key] = cleanMapping(value, updates)
      continue;
    }

    // Check reference
    if (value.$path === undefined)
      continue;

    let v = getValueFromStruct(value.$path, updates);
    if (v === undefined)
      delete mapping[key];
  }

  return mapping;
}

function transformMapping(mapping, updates, script) {
  let transformed = script(mapping);
  if (!transformed) {
    return null;
  }

  let mappings = [ transformed ];
  if (Array.isArray(transformed)) {
    mappings = transformed;
  }

  return mappings.map((mapping) => cleanMapping(mapping, updates));
}

function transformUpdates(updates, script) {
  let transformed = script(updates);
  if (!transformed) {
    return null;
  }

  let transformedUpdates = [ transformed ];
  if (Array.isArray(transformed)) {
    transformedUpdates = transformed;
  }

  return transformedUpdates;
}

function getResults(source, validSource, updates, transformedMappings, transformedUpdates) {

  return transformedMappings.map((transformedMapping, index) => { 

    let transformedUpdate = transformedUpdates[index];

    let results = {};

    for (let entry in validSource) {

      if (!entry.includes('.')) {
        continue;
      }

      let key = findTargetKey(entry, transformedMapping)
      if (key == null) {
        continue;
      }

      let value = getValueFromStruct(entry, updates);

      results[key] = value;
    }

    let mergedResults = Object.assign({}, transformedUpdate, results)
    let finalResults = trimInvalidKeys(mergedResults);

    // Append removed fields
    if (source['$removedFields']) {
      finalResults['$removedFields'] = source['$removedFields'];
    }

    return finalResults;
  });
}

function transform(source, script) {

  let validSource = trimInvalidKeys(source);
  let mapping = createMapping(validSource);
  let updates = getUpdates(validSource);
/*
  console.log('valid', validSource);
  console.log('mapping', mapping);
  console.log('updates', updates);
*/
  let transformedUpdates = transformUpdates(updates, script);
  if (!transformedUpdates) {
    return null;
  }

  let transformedMappings = transformMapping(mapping, updates, script);
  if (!transformedMappings) {
    return null;
  }
//  console.log('transformedMapping', transformedMappings);
//  console.log('transformedUpdates', transformedUpdates);

  let results = getResults(source, validSource, updates, transformedMappings, transformedUpdates);

  if (results.length === 1) {
    return results[0];
  }

  return results;

//  let transformed = script(updates);

  /*
  let results = {};

  for (let entry in validSource) {

    if (!entry.includes('.')) {
      continue;
    }

    let key = findTargetKey(entry, transformedMapping)
    if (key == null) {
      continue;
    }

    let value = getValueFromStruct(entry, updates);

    results[key] = value;
  }

  let mergedResults = Object.assign({}, transformed, results)
  let finalResults = trimInvalidKeys(mergedResults);

  // Append removed fields
  if (source['$removedFields']) {
    finalResults['$removedFields'] = source['$removedFields'];
  }

  return finalResults;
*/
}

//console.log(transform(source, script));

module.exports = {
  transform: transform,
};
