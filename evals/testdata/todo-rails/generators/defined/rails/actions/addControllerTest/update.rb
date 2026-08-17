  test "update applies what the body carried" do
    created = {{resource|model}}.create!(title: "before")
    put "/{{resource|table}}/#{created.id}", params: { title: "after", completed: true }, as: :json

    assert_response :success
    assert_equal "after", response.parsed_body["title"]
  end
