const form = document.getElementById("predictForm");
const fileInput = document.getElementById("fileInput");
const submitButton = document.getElementById("submitButton");
const previewContainer = document.getElementById("previewContainer");
const previewPlaceholder = document.getElementById("previewPlaceholder");
const previewImage = document.getElementById("previewImage");
const removeImageButton = document.getElementById("removeImageButton");
const statusElement = document.getElementById("status");
const resultElement = document.getElementById("result");
const breedValue = document.getElementById("breedValue");
const confidenceValue = document.getElementById("confidenceValue");
const filenameValue = document.getElementById("filenameValue");
const ALLOWED_IMAGE_TYPES = new Set(["image/jpeg", "image/png"]);

// Если задан window.CATBREED_API_URL в index.html, используем его.
// Иначе идем в локальный backend FastAPI.
const API_URL = window.CATBREED_API_URL || "http://127.0.0.1:8000/predict";

let currentPreviewUrl = "";
let selectedFile = null;

function setStatus(message, kind = "") {
    statusElement.textContent = message;
    statusElement.classList.remove("status--success", "status--error");
    if (kind === "success") {
        statusElement.classList.add("status--success");
    } else if (kind === "error") {
        statusElement.classList.add("status--error");
    }
}

function clearPreviewUrl() {
    if (currentPreviewUrl) {
        URL.revokeObjectURL(currentPreviewUrl);
        currentPreviewUrl = "";
    }
}

function resetPreview() {
    clearPreviewUrl();
    selectedFile = null;
    fileInput.value = "";
    previewImage.src = "";
    previewContainer.hidden = true;
    previewPlaceholder.hidden = false;
    resultElement.hidden = true;
    setStatus("");
    breedValue.textContent = "—";
    confidenceValue.textContent = "—";
    filenameValue.textContent = "—";
}

function showPreview(file) {
    clearPreviewUrl();

    currentPreviewUrl = URL.createObjectURL(file);
    previewImage.src = currentPreviewUrl;

    previewPlaceholder.hidden = true;
    previewContainer.hidden = false;
}

function formatConfidence(rawConfidence) {
    const numeric = Number(rawConfidence);
    if (!Number.isNaN(numeric)) {
        return `${(numeric * 100).toFixed(1)}%`;
    }
    return String(rawConfidence || "—");
}

fileInput.addEventListener("change", () => {
    const file = fileInput.files?.[0];

    resultElement.hidden = true;
    setStatus("");

    if (!file) {
        resetPreview();
        return;
    }

    if (!ALLOWED_IMAGE_TYPES.has(file.type)) {
        resetPreview();
        setStatus("Поддерживаются только JPG и PNG.", "error");
        return;
    }

    selectedFile = file;
    showPreview(file);
});

// Разрешаем выбрать тот же файл повторно: очищаем value перед открытием диалога.
fileInput.addEventListener("click", () => {
    fileInput.value = "";
});

removeImageButton.addEventListener("click", () => {
    resetPreview();
});

form.addEventListener("submit", async (event) => {
    event.preventDefault();

    const file = selectedFile || fileInput.files?.[0];

    if (!file) {
        setStatus("Сначала выберите изображение.", "error");
        resultElement.hidden = true;
        return;
    }

    submitButton.disabled = true;
    submitButton.textContent = "Определяем...";
    setStatus("Идёт обработка изображения...");
    resultElement.hidden = true;

    try {
        const formData = new FormData();
        formData.append("file", file);

        const response = await fetch(API_URL, {
            method: "POST",
            body: formData,
        });

        const contentType = response.headers.get("content-type") || "";
        let data;

        if (contentType.includes("application/json")) {
            data = await response.json();
        } else {
            const text = await response.text();
            throw new Error(`Ожидался JSON от API, но пришло: ${text.slice(0, 120)}`);
        }

        if (!response.ok) {
            const errorText = data?.detail || `Ошибка сервера: ${response.status}`;
            throw new Error(errorText);
        }

        breedValue.textContent = data.breed ?? "Неизвестно";
        confidenceValue.textContent = formatConfidence(data.confidence);
        filenameValue.textContent = data.filename ?? file.name;

        resultElement.hidden = false;
        setStatus("Результат успешно получен.", "success");
    } catch (error) {
        console.error(error);
        setStatus(`Не удалось получить результат: ${error.message}`, "error");
        resultElement.hidden = true;
    } finally {
        submitButton.disabled = false;
        submitButton.textContent = "Определить породу";
    }
});

window.addEventListener("beforeunload", () => {
    clearPreviewUrl();
});

window.addEventListener("paste", (event) => {
    const items = event.clipboardData?.items;
    if (!items || items.length === 0) {
        return;
    }

    for (const item of items) {
        if (item.kind !== "file") {
            continue;
        }

        const file = item.getAsFile();
        if (!file) {
            continue;
        }

        if (!ALLOWED_IMAGE_TYPES.has(file.type)) {
            setStatus("Из буфера можно вставлять только JPG или PNG.", "error");
            return;
        }

        selectedFile = new File([file], file.name || "pasted-image.png", { type: file.type });
        try {
            const dataTransfer = new DataTransfer();
            dataTransfer.items.add(selectedFile);
            fileInput.files = dataTransfer.files;
        } catch (_) {
            // В некоторых браузерах assignment fileInput.files может быть запрещен.
        }
        resultElement.hidden = true;
        showPreview(selectedFile);
        setStatus("Изображение вставлено из буфера. Нажми «Определить породу».", "success");
        return;
    }
});
