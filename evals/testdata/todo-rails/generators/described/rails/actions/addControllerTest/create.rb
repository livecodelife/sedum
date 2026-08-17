  test "create returns the todo it stored" do
    post "/{{resource|table}}", params: { title: "write it down" }, as: :json

    assert_response :created
    assert_equal "write it down", response.parsed_body["title"]
  end
