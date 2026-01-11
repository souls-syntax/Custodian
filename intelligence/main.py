import torch
import torch.nn.functional as F
from fastapi import FastAPI
from pydantic import BaseModel
from transformers import DistilBertTokenizer, DistilBertForSequenceClassification
from ddgs import DDGS
from google import genai
import random
import os
import time
import json

app = FastAPI()

# ==========================================
# CONFIGURATION
# ==========================================
API_KEY = os.getenv("GOOGLE_API_KEY") or os.getenv("GEMINI_API_KEY")
if not API_KEY:
    print("⚠️ CRITICAL: NO API KEY FOUND.")
    client = None
else:
    client = genai.Client(api_key=API_KEY)

MODEL_PATH = "./models/template_model_kohai_v1"
device = torch.device("cpu")

# Load BERT
try:
    model = DistilBertForSequenceClassification.from_pretrained(MODEL_PATH)
    tokenizer = DistilBertTokenizer.from_pretrained(MODEL_PATH)
    model.to(device)
    model.eval()
    MODEL_LOADED = True
    print("✅ BERT ONLINE.")
except:
    MODEL_LOADED = False
    print("⚠️ BERT OFFLINE.")

class QueryRequest(BaseModel):
    query_text: str

# ==========================================
# 1. HELPER: SMART SEARCH
# ==========================================
def smart_search(user_query):
    # STEP 1: Ask Gemini to generate a Google-friendly search term
    try:
        extraction_prompt = f"Convert this user claim into a simple 5-7 word Google search query. Output ONLY the search query. Claim: {user_query}"
        
        search_term_resp = client.models.generate_content(
            model="gemini-2.5-flash", 
            contents=extraction_prompt
        )
        search_term = search_term_resp.text.strip()
        print(f"🧠 Optimized Search Term: '{search_term}'")
    except:
        search_term = user_query # Fallback

    # STEP 2: DuckDuckGo Search
    try:
        print(f"🔎 DDGS Radar Active: Scanning for '{search_term}'")
        # limit to 3 results to keep it fast
        with DDGS() as ddgs:
            # .text() returns a generator, so we wrap it in list()
            results = list(ddgs.text(search_term, max_results=5))
        
        context = ""
        count = 0
        for i, res in enumerate(results):
            # Safe .get() to avoid crashing on missing keys
            body = res.get('body', res.get('snippet', ''))
            href = res.get('href', '#')
            context += f"Source {i+1}: {body} (Link: {href})\n"
            count += 1

        
        print(f"Found {count} sources.")
        return context if count > 0 else "No relevant news found."
    except Exception as e:
        print(f"❌ Search Failed: {e}")
        return "Search Tool Unavailable."

# ==========================================
# 2. ENDPOINTS
# ==========================================
@app.post("/predict")
async def predict(req: QueryRequest):
    if "?" in req.query_text or len(req.query_text.split()) < 10:
        return {"verdict": "uncertain", "confidence": 0.45, "decider": "Guardrail-Routing"}

    if not MODEL_LOADED:
        return {"verdict": "error", "confidence": 0.0}

    inputs = tokenizer(req.query_text, return_tensors="pt", truncation=True, padding=True, max_length=512).to(device)
    with torch.no_grad():
        outputs = model(**inputs)
    probs = F.softmax(outputs.logits, dim=-1)
    conf, class_id = torch.max(probs, dim=1)
    labels = {0: "likely_true", 1: "likely_false"}
    
    return {
        "verdict": labels[class_id.item()],
        "confidence": float(conf.item()),
        "decider": "DistilBERT-v1"
    }


SIMULATION_MODE = False
@app.post("/predict_llm")
def predict_llm(req: QueryRequest):
    if SIMULATION_MODE:
        print(f"simulating {req.query_text}")
        time.sleep(random.uniform(1.0, 3.0)) 
        return {
            "verdict": random.choice(["likely_true", "likely_false", "uncertain"]),
            "confidence": random.uniform(0.7, 0.99),
            "decider": "Gemini-Agent (SIMULATION)",
            "explanation": "Simulated response for stress testing."
        }

    query = req.query_text
    print(f"🤖 AGENT ACTIVATED: Processing '{query}'")

    if not client:
        return {"verdict": "error", "confidence": 0.0}

    # 1. SMART SEARCH
    context_data = smart_search(query)
    
    # 2. VERIFICATION PROMPT
    prompt = f"""
    You are CUSTODIAN, a fact-checking AI.
    
    LIVE EVIDENCE:
    {context_data}

    USER CLAIM:
    {query}

    INSTRUCTIONS:
    - Use the LIVE EVIDENCE to judge the claim.
    - If the evidence confirms the event (even if it's recent), verdict is 'likely_true'.
    - If evidence contradicts, 'likely_false'.
    - If evidence is empty, 'uncertain'.
    - Return confidence between 0 and 1 float
    
    Return JSON: {{ "verdict": "...", "confidence": float, "explanation": "..." }}
    """

    try:
        response = client.models.generate_content(
            model="gemini-2.5-flash",
            contents=prompt
        )
        # Parse JSON
        txt = response.text.replace("```json", "").replace("```", "").strip()
        res = json.loads(txt)
        return {
            "verdict": res.get("verdict", "uncertain"),
            "confidence": res.get("confidence", 0.5),
            "decider": "Gemini-Agent",
            "explanation": res.get("explanation")
        }
    except Exception as e:
        return {"verdict": "error", "decider": f"Error: {e}"}

@app.get("/healthz")
def healthz():
    return {
            "Health": "ok"
            }
