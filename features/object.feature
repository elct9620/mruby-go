Feature: Object
  Scenario: I can new a object
    When I execute ruby code:
      """
      Object.new
      """
    Then there should return object

  Scenario: I can call a method defined in the object
    When I execute ruby code:
      """
      class Hello
        def world
          "world"
        end
      end
      Hello.new.world
      """
    Then there should return string "world"

  Scenario: I can new a object
    When I execute ruby code:
      """
      Object.inspect
      """
    Then there should return string "Object"
