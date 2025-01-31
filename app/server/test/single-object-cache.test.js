import { SingleObjectCache } from '/assets/single-object-cache.js'

QUnit.module('SingleObjectCache', function () {
  QUnit.test('simple set get', function (assert) {
    const cache = new SingleObjectCache()
    const expectedValue = 'value'
    const key = { bool: true, string: 'string', number: 1 }

    cache.set(key, expectedValue)

    assert.true(cache.has(key))

    assert.equal(cache.get(key), expectedValue)
  })

  QUnit.test(
    'get different object with equal values and same order',
    function (assert) {
      const cache = new SingleObjectCache()
      const expectedValue = 'value'
      const key = { bool: true, string: 'string', number: 1 }

      cache.set(key, expectedValue)

      const newKey = { bool: true, string: 'string', number: 1 }

      assert.true(cache.has(newKey))

      assert.equal(cache.get(newKey), expectedValue)
    }
  )

  QUnit.test(
    'get different object with equal values and different order',
    function (assert) {
      const cache = new SingleObjectCache()
      const expectedValue = 'value'
      const key = { bool: true, string: 'string', number: 1 }

      cache.set(key, expectedValue)

      const newKey = { string: 'string', bool: true, number: 1 }

      assert.true(cache.has(newKey))

      assert.equal(cache.get(newKey), expectedValue)
    }
  )

  QUnit.test('get different object with different values', function (assert) {
    const cache = new SingleObjectCache()
    const expectedValue = 'value'
    const key = { bool: true, string: 'string', number: 1 }

    cache.set(key, expectedValue)

    const newKey = { string: 'string2', bool: true, number: 1 }

    assert.false(cache.has(newKey))

    assert.notEqual(cache.get(newKey), expectedValue)
    assert.equal(cache.get(newKey), null)
  })

  QUnit.test('set new cache value', function (assert) {
    const cache = new SingleObjectCache()
    const initialValue = 'initialValue'
    const expectedValue = 'value'
    const key = { bool: true, string: 'string', number: 1 }

    cache.set(key, initialValue)

    const newKey = { string: 'string2', bool: true, number: 1 }

    cache.set(newKey, expectedValue)

    assert.false(cache.has(key))
    assert.notEqual(cache.get(key), initialValue)
    assert.equal(cache.get(key), null)

    assert.true(cache.has(newKey))
    assert.equal(cache.get(newKey), expectedValue)
  })
})
