package go_cake

import "context"

type DatabaseDriver interface {
	GetUnderlyingDriver() any

	TestModel(
		idField string,
		etagField string,
		model GoCakeModel,
		dbPath string) error

	Find(
		model GoCakeModel,
		where, sort string,
		page, perPage int64,
		ctx context.Context,
		userData any) ([]GoCakeModel, HTTPError)

	Delete(
		model GoCakeModel,
		documents []GoCakeModel,
		ctx context.Context,
		userData any) HTTPError

	Total(
		model GoCakeModel,
		where string,
		ctx context.Context,
		userData any) (uint64, HTTPError)

	Insert(
		model GoCakeModel,
		documents []GoCakeModel,
		ctx context.Context,
		userData any) HTTPError

	Update(
		model GoCakeModel,
		documents []GoCakeModel,
		ctx context.Context,
		userData any) HTTPError

	GetWhereFields(model GoCakeModel, where string) ([]string, HTTPError)
	GetSortFields(model GoCakeModel, sort string) ([]string, HTTPError)
}
