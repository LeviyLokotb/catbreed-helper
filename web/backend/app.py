from fastapi import FastAPI, UploadFile, File
from fastapi.middleware.cors import CORSMiddleware

app = FastAPI(title="Cat Breed Helper API")

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

@app.get("/")
def root():
    return {"ok": True, "message": "Cat Breed Helper API is running"}

@app.post("/predict")
async def predict(file: UploadFile = File(...)):
    return {
        "success": True,
        "breed": "Британская короткошёрстная",
        "breed_id": "british",
        "confidence": 0.95,
        "source": "mock"
    }