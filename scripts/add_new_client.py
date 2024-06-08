import os

from dotenv import load_dotenv
from supabase import Client, create_client

load_dotenv()

# Make this a parameter
clientName = "flock3" 

supabase: Client = create_client(os.environ.get("SUPABASE_URL"), os.environ.get("SUPABASE_KEY"))

data, count = supabase.table("clients").insert({"name": "flock3"}).execute()

print(data, count)
