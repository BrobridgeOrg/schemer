const assert = require('assert');
const { transform } = require('../js/core');

describe('Transformer testing', function() {

  let source = {
    '$removedFields': [ 'unused_field' ],
    'aaa.bbb.ccc': 123,
    'aaa.bbb.sss': 'SSS',
    'aaa.ddd': 456,
    'aaa.ggg': {
      ooo: 666,
      nested: {
        deep: 'deep',
        nested: {
          deeper: 'deeper'
        }
      }
    },
    'arr.list.1': 'new_element',
    'arr1.list.2': 'new_element',
    'arr1.map.0': {
      id: 1,
      name: 'test_1',
    },
    'arr1.map.1': {
      id: 2,
      name: 'test_2',
    },
    'arr1.map.2.name': 'test_3',
    'aa': 'short',
    'xxx': 999,
    'myArray': [
      'arr1',
      'arr2'
    ],
    'object': {
      title: 'title',
    }
  };

  function script(source) {
    return {
      'zzz': source.aaa,
      'kkk': source.aaa.ddd + 1,
      'qqq': source.aaa.ggg,
      'gggzzz': source.aaa.ggg.zzz, // empty
      'jjj': source.xxx,
      'new_aa': source.aa,
      'arrA': source.arr1,
      'nested': {
        sss: source.aaa.bbb.sss,
        ooo: source.aaa.ggg.ooo,
        deep: source.aaa.ggg.nested.deep,
        deeper: source.aaa.ggg.nested.nested.deeper
      },
      'myArray': source.myArray,
      'myMap': source.arr1.map[0],
      'title': source.object.title
    };
  }

  it('$removedFields', function() {

    let result = transform(source, script);

    assert.deepEqual(result.$removedFields, [ 'unused_field' ]);
  });

  it('Update specific fields', function() {

    let result = transform(source, script);

    // transformed field
    assert.equal(result['zzz.bbb.ccc'], source['aaa.bbb.ccc']);
    assert.equal(result['zzz.bbb.sss'], source['aaa.bbb.sss']);
    assert.equal(result['zzz.ddd'], source['aaa.ddd']);
    assert.deepEqual(result['zzz.ggg'], source['aaa.ggg']);

    // array and map
    assert.equal(result['arrA.list.2'], source['arr1.list.2']);
    assert.deepEqual(result['arrA.map.0'], source['arr1.map.0']);
    assert.deepEqual(result['arrA.map.1'], source['arr1.map.1']);
    assert.equal(result['arrA.map.2.name'], source['arr1.map.2.name']);
    
    // map
    assert.deepEqual(result.myArray, source.myArray);
    assert.deepEqual(result.myMap, source['arr1.map.0']);

  });

  it('Nested structure', function() {

    let result = transform(source, script);

    assert.equal(result.nested.sss, source['aaa.bbb.sss']);
    assert.equal(result.nested.ooo, source['aaa.ggg'].ooo);
    assert.equal(result.nested.deep, source['aaa.ggg'].nested.deep);
    assert.equal(result.nested.deeper, source['aaa.ggg'].nested.nested.deeper);
  });
});
