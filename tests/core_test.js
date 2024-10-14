
const assert = require('assert');
const { execute } = require('../js/core');

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
      meta: {
        description: 'description',
        keywords: [
          'keyword1',
          'keyword2'
        ]
      }
    }
  };

  it('direct', function() {

    function script(source = {}) {
      return source;
    }

    let result = execute(source, script);

    for (let field in result) {
      assert.deepEqual(result[field], source[field]);
    }
  });

  it('mapping', function() {

    function script(source = {}) {

      return source.$mapping({
        AAA: 'aaa',
        BBB: 'aaa.bbb',
        XXX: 'xxx'
      });
    }

    let result = execute(source, script);

    assert.deepEqual(result['AAA'], source['aaa']);
    assert.deepEqual(result['AAA.ddd'], source['aaa.ddd']);
    assert.deepEqual(result['AAA.ggg'], source['aaa.ggg']);
    assert.deepEqual(result['BBB.ccc'], source['aaa.bbb.ccc']);
    assert.deepEqual(result['BBB.sss'], source['aaa.bbb.sss']);
    assert.deepEqual(result['XXX'], source['xxx']);
  });

  it('$removedFields', function() {

    function script(source = {}) {

      return source.$mapping({
        UNUSED: 'unused_field',
      })
    }

    let result = execute(source, script);

    assert.deepEqual(result.$removedFields, [ 'UNUSED' ]);
  });

  it('Check Visibility', function() {

    function script(source = {}) {

      return source.$mapping({
        UNUSED: 'unused_field',
        AAA: 'aaa',
        BBB: 'aaa.bbb',
        XXX: 'xxx'
      })
    }

    let result = execute(source, script);

    for (let field in result) {
      if (field == '$removedFields')
        continue;

      if (field.startsWith('$')) {
        assert.fail('Unexpected field: ' + field);
      }
    }
  });

  it('Update specific fields', function() {

    function script(source = {}) {

      return source.$mapping({
        zzz: 'aaa',
        arrA: 'arr1',
        myArray: 'myArray',
        myMap: 'arr1.map.0',
        empty: 'source.empty',
      });
    }

    let result = execute(source, script);

    //console.log(result);

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

    // empty
    assert.equal(result.empty, undefined);
  });

  it('Watch specific fields', function() {

    function script(source = {}) {

      let output = {};

      source.$watch('aaa.bbb.ccc', (value) => {
        output.AAA = value;
      });

      source.$watch('aaa.bbb.sss', (value) => {
        output.SSS = value;
      });

      source.$watch('unused_field', (value, removed) => {
        if (!output.$removedFields)
          output.$removedFields = [];

        output.$removedFields.push('UNUSED');
      });

      return output;
    }

    let result = execute(source, script);
    
    assert.equal(result['AAA'], source['aaa.bbb.ccc']);
    assert.equal(result['SSS'], source['aaa.bbb.sss']);
    assert.deepEqual(result.$removedFields, [ 'UNUSED' ]);
  });

  it('With operators', function() {

    function script(source = {}) {

      return {
          kkk: source.$getValue('aaa.ddd') + 99999,
          op_str: source.$getValue('aaa.ddd') + 'QQQQQ',
          op_empty: source.$getValue('aaa.empty')
        };
    }

    let result = execute(source, script);
    
    assert.equal(result['kkk'], source['aaa.ddd'] + 99999);
    assert.equal(result['op_str'], source['aaa.ddd'] + 'QQQQQ');
    assert.equal(result['op_empty'], undefined);
  });

  it('Nested structure', function() {

    function script(source = {}) {

      return {
        'nested': {
          sss: source.$getValue('aaa.bbb.sss'),
          ooo: source.$getValue('aaa.ggg.ooo'),
          deep: source.$getValue('aaa.ggg.nested.deep'),
          deeper: source.$getValue('aaa.ggg.nested.nested.deeper'),
        },
        'title': source.$getValue('object.title'),
        'desc': source.$getValue('object.meta.description'),
        'metadata': {
          keywords: source.$getValue('object.meta.keywords')
        }
      };
    }

    let result = execute(source, script);

    // From specified fields
    assert.equal(result.nested.sss, source['aaa.bbb.sss']);
    assert.equal(result.nested.ooo, source['aaa.ggg'].ooo);
    assert.equal(result.nested.deep, source['aaa.ggg'].nested.deep);
    assert.equal(result.nested.deeper, source['aaa.ggg'].nested.nested.deeper);

    // From static object
    assert.equal(result.title, source.object.title);
    assert.equal(result.desc, source.object.meta.description);
    assert.deepEqual(result.metadata.keywords, source.object.meta.keywords);
  });

  it('Access array', function() {

    function script(source = {}) {

      return {
        'newArr': source.$getValue('myArray'),
        'firstElement': source.$getValue('myArray.0'),
        'secondElement': source.$getValue('myArray.1'),
      };
    }

    let result = execute(source, script);

    assert.deepEqual(result.newArr, source.myArray);
    assert.equal(result['firstElement'], source.myArray[0]);
    assert.equal(result['secondElement'], source.myArray[1]);
  });
});
