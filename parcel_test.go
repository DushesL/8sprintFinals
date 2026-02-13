package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	// randSource источник псевдо случайных чисел.
	// Для повышения уникальности в качестве seed
	// используется текущее время в unix формате (в виде числа)
	randSource = rand.NewSource(time.Now().UnixNano())
	// randRange использует randSource для генерации случайных чисел
	randRange = rand.New(randSource)
)

// getTestParcel возвращает тестовую посылку
func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE parcel (
			number INTEGER PRIMARY KEY AUTOINCREMENT,
			client INTEGER NOT NULL,
			status TEXT NOT NULL,
			address TEXT NOT NULL,
		created_at TEXT NOT NULL
	)`)
	require.NoError(t, err)

	store := NewParcelStore(db)

	// add
	parcel := getTestParcel()
	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.Greater(t, id, 0)
	parcel.Number = id

	// get
	got, err := store.Get(id)
	require.NoError(t, err)
	assert.Equal(t, parcel, got)

	// delete
	err = store.Delete(id)
	require.NoError(t, err)

	// проверка, что посылку больше нельзя получить
	_, err = store.Get(id)
	require.Error(t, err)
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE parcel (
			number INTEGER PRIMARY KEY AUTOINCREMENT,
			client INTEGER NOT NULL,
			status TEXT NOT NULL,
			address TEXT NOT NULL,
		created_at TEXT NOT NULL
	)`)
	require.NoError(t, err)

	store := NewParcelStore(db)

	// add
	parcel := getTestParcel()
	id, err := store.Add(parcel)
	require.NoError(t, err)
	parcel.Number = id

	// set address
	newAddress := "new test address"
	err = store.SetAddress(id, newAddress)
	require.NoError(t, err)

	// check
	got, err := store.Get(id)
	require.NoError(t, err)
	assert.Equal(t, newAddress, got.Address)
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE parcel (
			number INTEGER PRIMARY KEY AUTOINCREMENT,
			client INTEGER NOT NULL,
			status TEXT NOT NULL,
			address TEXT NOT NULL,
		created_at TEXT NOT NULL
	)`)
	require.NoError(t, err)

	store := NewParcelStore(db)

	// add
	parcel := getTestParcel()
	id, err := store.Add(parcel)
	require.NoError(t, err)
	parcel.Number = id

	// set status
	err = store.SetStatus(id, ParcelStatusSent)
	require.NoError(t, err)

	// check
	got, err := store.Get(id)
	require.NoError(t, err)
	assert.Equal(t, ParcelStatusSent, got.Status)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE parcel (
			number INTEGER PRIMARY KEY AUTOINCREMENT,
			client INTEGER NOT NULL,
			status TEXT NOT NULL,
			address TEXT NOT NULL,
		created_at TEXT NOT NULL
	)`)
	require.NoError(t, err)

	store := NewParcelStore(db)

	parcels := []Parcel{getTestParcel(), getTestParcel(), getTestParcel()}
	parcelMap := make(map[int]Parcel)

	client := randRange.Intn(10_000_000)
	for i := range parcels {
		parcels[i].Client = client
	}

	// add
	for _, p := range parcels {
		id, err := store.Add(p)
		require.NoError(t, err)
		p.Number = id
		parcelMap[id] = p
	}

	// get by client
	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err)
	assert.Len(t, storedParcels, len(parcels))

	// check
	for _, got := range storedParcels {
		expected, ok := parcelMap[got.Number]
		assert.True(t, ok)
		assert.Equal(t, expected.Client, got.Client)
		assert.Equal(t, expected.Status, got.Status)
		assert.Equal(t, expected.Address, got.Address)
		assert.Equal(t, expected.CreatedAt, got.CreatedAt)
	}
}
