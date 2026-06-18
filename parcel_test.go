package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	randSource = rand.NewSource(time.Now().UnixNano())
	randRange  = rand.New(randSource)
)

func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

func TestAddGetDelete(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)

	parcel := getTestParcel()

	number, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotZero(t, number)

	parcel.Number = number

	got, err := store.Get(number)
	require.NoError(t, err)
	require.Equal(t, parcel, got)

	err = store.Delete(number)
	require.NoError(t, err)

	_, err = store.Get(number)
	require.Error(t, err)
}

func TestSetAddress(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)

	parcel := getTestParcel()

	number, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotZero(t, number)

	newAddress := "new test address"

	err = store.SetAddress(number, newAddress)
	require.NoError(t, err)

	ch, err := store.Get(number)
	require.NoError(t, err)
	require.Equal(t, newAddress, ch.Address)
}

func TestSetStatus(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)

	parcel := getTestParcel()

	number, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotZero(t, number)

	status := ParcelStatusSent
	err = store.SetStatus(number, status)
	require.NoError(t, err)

	ch, err := store.Get(number)
	require.NoError(t, err)
	require.Equal(t, status, ch.Status)
}

func TestGetByClient(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}

	parcelMap := map[int]Parcel{}

	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	for i := 0; i < len(parcels); i++ {
		id, err := store.Add(parcels[i])
		require.NoError(t, err)
		require.NotZero(t, id)

		parcels[i].Number = id

		parcelMap[id] = parcels[i]
	}

	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err)
	require.Equal(t, len(parcels), len(storedParcels))

	for _, parcel := range storedParcels {
		value, ok := parcelMap[parcel.Number]
		require.True(t, ok)
		require.Equal(t, value, parcel)
	}
}
