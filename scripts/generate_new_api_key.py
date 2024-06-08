import os

from dotenv import load_dotenv
from supabase import Client, create_client

load_dotenv()

# Make this a parameter
clientName = "flocasdk1"

supabase: Client = create_client(
    os.environ.get("SUPABASE_URL"), os.environ.get("SUPABASE_KEY")
)

data, count = supabase.table("clients").select("*").eq("name", clientName).execute()

print(len(data[1]))

print(data[1][0]['created_at'])