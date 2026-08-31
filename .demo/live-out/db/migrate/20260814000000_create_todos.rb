class CreateTodos < ActiveRecord::Migration[7.2]
  def change
    create_table :todos do |t|
      # sedum:anchor:columns
      # sedum:addColumn {"tier":"owned","record":"eval-todo-rails","kwargs":{"default":"nil","name":"title","nullable":false,"resource":"todo","stamp":"20260814000000","type":"string"}}
      t.string :title, null: false, default: nil
      # /sedum:addColumn
      # sedum:addColumn {"tier":"owned","record":"eval-todo-rails","kwargs":{"default":"false","name":"completed","nullable":false,"resource":"todo","stamp":"20260814000000","type":"boolean"}}
      t.boolean :completed, null: false, default: false
      # /sedum:addColumn

      t.timestamps
    end
  end
end
