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

  Scenario: I can inspect object which is an instance of a class
    When I execute ruby code:
      """
      class Hello
        def world
          "world"
        end
      end
      Hello.new.inspect
      """
    Then there should return string "Hello"
