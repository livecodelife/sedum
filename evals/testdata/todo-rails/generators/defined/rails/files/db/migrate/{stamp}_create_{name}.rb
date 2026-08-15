class Create{{name|models}} < ActiveRecord::Migration[7.2]
  def change
    create_table :{{name|table}} do |t|
      # sedum:anchor:columns

      t.timestamps
    end
  end
end
