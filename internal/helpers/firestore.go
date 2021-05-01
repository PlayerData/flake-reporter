package helpers

import (
	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

func FirestoreForEachDocument(iter *firestore.DocumentRefIterator, f func(item *firestore.DocumentRef) error) error {
	for {
		item, err := iter.Next()

		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}

		err = f(item)
		if err != nil {
			return err
		}
	}

	return nil
}

func FirestoreForEachCollection(iter *firestore.CollectionIterator, f func(item *firestore.CollectionRef) error) error {
	for {
		item, err := iter.Next()

		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}

		err = f(item)
		if err != nil {
			return err
		}
	}

	return nil
}
