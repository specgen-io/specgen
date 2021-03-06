package gents

import (
	"github.com/specgen-io/specgen/v2/gen"
	"strings"
)

var IoTs = "io-ts"

var importIoTsEncoding = `import * as t from './io-ts'`

func generateIoTsStaticCode(path string) *gen.TextFile {
	code := `
export * from 'io-ts'
export * from 'io-ts-types'

import * as t from 'io-ts'

import { pipe, identity } from 'fp-ts/lib/function'
import { fold, map } from 'fp-ts/lib/Either'


enum Enum {}

export class EnumType<E extends typeof Enum> extends t.Type<E[keyof E]> {
  readonly _tag: 'EnumType' = 'EnumType'
  private readonly _enum: E
  private readonly _enumValues: Set<string | number>
  constructor(e: E, name?: string) {
    super(
      name || 'enum',
      (u): u is E[keyof E] => {
        if (!this._enumValues.has(u as any)) return false
        // Don't allow key names from number enum reverse mapping
        if (typeof (this._enum as any)[u as string] === 'number') return false
        return true
      },
      (u, c) => (this.is(u) ? t.success(u) : t.failure(u, c)),
      t.identity
    )
    this._enum = e
    this._enumValues = new Set(Object.values(e))
  }
}

const enumType = <E extends typeof Enum>(e: E, name?: string) => new EnumType<E>(e, name)

export { enumType as enum }

export class WithDefault<RT extends t.Any, A = any, O = A, I = unknown> extends t.Type<A, O, I> {
  readonly _tag: 'WithDefault' = 'WithDefault'
  constructor(
      name: string,
      is: WithDefault<RT, A, O, I>['is'],
      validate: WithDefault<RT, A, O, I>['validate'],
      serialize: WithDefault<RT, A, O, I>['encode'],
      readonly type: RT,
  ) {
      super(name, is, validate, serialize)
  }
}

export const withDefault = <RT extends t.Type<A, O>, A = any, O = A>(type: RT, defaultValue: t.TypeOf<RT>): WithDefault<RT, t.TypeOf<RT>, t.OutputOf<RT>, unknown> => {
  const Nullable = t.union([type, t.null, t.undefined])
  return new WithDefault(
    'WithDefault',
      (m: unknown): m is t.TypeOf<RT> => type.is(m),
      (s: unknown, c: t.Context): t.Validation<t.TypeOf<RT>> => {
          const validationResult: t.Validation<t.TypeOf<RT | t.NullC | t.UndefinedC>> = Nullable.validate(s, c)
          const applyDefault = map<A | null | undefined, A>(value => value != null ? value : defaultValue)
          return applyDefault(validationResult)
      },
      (a: t.TypeOf<RT>) => type.encode(a),
      type,
  )
}

export class DecodeError extends Error {
    errors: t.Errors
    constructor(errors: t.Errors) {
        super('Decoding failed')
        this.errors = errors
    }
}

export const decode = <A, O, I>(codec: t.Type<A, O, I>, value: I): A => {
    return pipe(
        codec.decode(value),
        fold(
            errors => { throw new DecodeError(errors) },
            identity
        )
    )
}

export const encode = <A, O, I>(codec: t.Type<A, O, I>, value: A): O => {
    return codec.encode(value)
}
`
	return &gen.TextFile{path, strings.TrimSpace(code)}
}
