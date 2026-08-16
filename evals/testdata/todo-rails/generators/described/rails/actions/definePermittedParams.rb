  # Permitting on the root parameters rather than requiring the resource key
  # first means a bare JSON body works whether or not parameter wrapping is
  # configured, and an empty body yields an empty hash instead of raising.
  #
  # The attribute list is split rather than interpolated as Ruby so that the
  # caller writes names and nothing else - no colons, no commas to get right.
  def {{resource|record}}_params
    params.permit(*"{{attributes}}".split(/[,\s]+/).reject(&:empty?).map(&:to_sym))
  end
