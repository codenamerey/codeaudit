first_name = "Elon"
last_name = "Musk"


def get_full_name(first_name, last_name):
  return first_name + last_name

def get_first_name(name):
  return name

def greetPerson(name, age):
  print(f"Hello {name}! How are you doing today?")
  if age < 18:
    print("Where are your parents and shouldn't you be in school?")
  else:
      print("Would you like to invest in my new company?")
