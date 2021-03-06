package tests

import (
	"fmt"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

// Insert a document into the table posts using a struct.
func ExampleTerm_Insert_struct() {
	type Post struct {
		ID      int    `rethinkdb:"id"`
		Title   string `rethinkdb:"title"`
		Content string `rethinkdb:"content"`
	}

	resp, err := r.DB("examples").Table("posts").Insert(Post{
		ID:      1,
		Title:   "Lorem ipsum",
		Content: "Dolor sit amet",
	}).RunWrite(session)
	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Printf("%d row inserted", resp.Inserted)

	// Output:
	// 1 row inserted
}

// Insert a document without a defined primary key into the table posts where
// the primary key is id.
func ExampleTerm_Insert_generatedKey() {
	type Post struct {
		Title   string `rethinkdb:"title"`
		Content string `rethinkdb:"content"`
	}

	resp, err := r.DB("examples").Table("posts").Insert(map[string]interface{}{
		"title":   "Lorem ipsum",
		"content": "Dolor sit amet",
	}).RunWrite(session)
	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Printf("%d row inserted, %d key generated", resp.Inserted, len(resp.GeneratedKeys))

	// Output:
	// 1 row inserted, 1 key generated
}

// Insert a document into the table posts using a map.
func ExampleTerm_Insert_map() {
	resp, err := r.DB("examples").Table("posts").Insert(map[string]interface{}{
		"id":      2,
		"title":   "Lorem ipsum",
		"content": "Dolor sit amet",
	}).RunWrite(session)
	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Printf("%d row inserted", resp.Inserted)

	// Output:
	// 1 row inserted
}

// Insert multiple documents into the table posts.
func ExampleTerm_Insert_multiple() {
	resp, err := r.DB("examples").Table("posts").Insert([]interface{}{
		map[string]interface{}{
			"title":   "Lorem ipsum",
			"content": "Dolor sit amet",
		},
		map[string]interface{}{
			"title":   "Lorem ipsum",
			"content": "Dolor sit amet",
		},
	}).RunWrite(session)
	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Printf("%d rows inserted", resp.Inserted)

	// Output:
	// 2 rows inserted
}

// Insert a document into the table posts, replacing the document if it already
// exists.
func ExampleTerm_Insert_upsert() {
	resp, err := r.DB("examples").Table("posts").Insert(map[string]interface{}{
		"id":    1,
		"title": "Lorem ipsum 2",
	}, r.InsertOpts{
		Conflict: "replace",
	}).RunWrite(session)
	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Printf("%d row replaced", resp.Replaced)

	// Output:
	// 1 row replaced
}

// Update the status of the post with id of 1 to published.
func ExampleTerm_Update() {
	resp, err := r.DB("examples").Table("posts").Get(2).Update(map[string]interface{}{
		"status": "published",
	}).RunWrite(session)
	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Printf("%d row replaced", resp.Replaced)

	// Output:
	// 1 row replaced
}

// Update bob's cell phone number.
func ExampleTerm_Update_nested() {
	resp, err := r.DB("examples").Table("users").Get("bob").Update(map[string]interface{}{
		"contact": map[string]interface{}{
			"phone": "408-555-4242",
		},
	}).RunWrite(session)
	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Printf("%d row replaced", resp.Replaced)

	// Output:
	// 1 row replaced
}

// Update the status of all posts to published.
func ExampleTerm_Update_all() {
	resp, err := r.DB("examples").Table("posts").Update(map[string]interface{}{
		"status": "published",
	}).RunWrite(session)
	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Printf("%d row replaced", resp.Replaced)

	// Output:
	// 4 row replaced
}

// Increment the field view of the post with id of 1. If the field views does not
// exist, it will be set to 0.
func ExampleTerm_Update_increment() {
	resp, err := r.DB("examples").Table("posts").Get(1).Update(map[string]interface{}{
		"views": r.Row.Field("views").Add(1).Default(0),
	}).RunWrite(session)
	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Printf("%d row replaced", resp.Replaced)

	// Output:
	// 1 row replaced
}

// Update the status of the post with id of 1 using soft durability.
func ExampleTerm_Update_softDurability() {
	resp, err := r.DB("examples").Table("posts").Get(2).Update(map[string]interface{}{
		"status": "draft",
	}, r.UpdateOpts{
		Durability: "soft",
	}).RunWrite(session)
	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Printf("%d row replaced", resp.Replaced)

	// Output:
	// 1 row replaced
}

// Delete a single document from the table posts.
func ExampleTerm_Delete() {
	resp, err := r.DB("examples").Table("posts").Get(2).Delete().RunWrite(session)
	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Printf("%d row deleted", resp.Deleted)

	// Output:
	// 1 row deleted
}

// Delete all comments where the field status is published
func ExampleTerm_Delete_many() {
	resp, err := r.DB("examples").Table("posts").Filter(map[string]interface{}{
		"status": "published",
	}).Delete().RunWrite(session)
	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Printf("%d rows deleted", resp.Deleted)

	// Output:
	// 4 rows deleted
}

func ExampleTerm_SetWriteHook() {
	resp, err := r.DB("test").Table("test").SetWriteHook(
		func(id r.Term, oldVal r.Term, newVal r.Term) r.Term {
			return r.Branch(oldVal.And(newVal),
				newVal.Merge(map[string]r.Term{"write_counter": oldVal.Field("write_counter").Add(1)}),
				newVal,
				newVal.Merge(r.Expr(map[string]int{"write_counter": 1})),
				nil,
			)
		}).RunWrite(session)

	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Printf("%d hook created", resp.Created)
	// Output:
	// 1 hook created
}
