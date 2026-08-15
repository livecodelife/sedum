  test "a {{resource|record}} without {{attribute}} is invalid" do
    {{resource|record}} = {{resource|model}}.new

    assert_not {{resource|record}}.valid?
    assert_includes {{resource|record}}.errors.attribute_names, :{{attribute}}
  end
