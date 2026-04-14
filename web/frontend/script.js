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
const modeValue = document.getElementById("modeValue");

const API_URL = "https://catbreed-helper.onrender.com/predict";

let currentPreviewUrl = "";

function clearPreviewUrl() {
    if (currentPreviewUrl) {
        URL.revokeObjectURL(currentPreviewUrl);
        currentPreviewUrl = "";
    }
}

function resetPreview() {
    clearPreviewUrl();
    fileInput.value = "";
    previewImage.src = "";
    previewContainer.hidden = true;
    previewPlaceholder.hidden = false;
    resultElement.hidden = true;
    statusElement.textContent = "";
    breedValue.textContent = "—";
    confidenceValue.textContent = "—";
    modeValue.textContent = "—";
}

function showPreview(file) {
    clearPreviewUrl();

    currentPreviewUrl = URL.createObjectURL(file);
    previewImage.src = currentPreviewUrl;

    previewPlaceholder.hidden = true;
    previewContainer.hidden = false;
}

fileInput.addEventListener("change", () => {
    const file = fileInput.files?.[0];

    resultElement.hidden = true;
    statusElement.textContent = "";

    if (!file) {
        resetPreview();
        return;
    }

    if (!file.type.startsWith("image/")) {
        resetPreview();
        statusElement.textContent = "Пожалуйста, выберите именно изображение.";
        return;
    }

    showPreview(file);
});

removeImageButton.addEventListener("click", () => {
    resetPreview();
});

form.addEventListener("submit", async (event) => {
    event.preventDefault();

    const file = fileInput.files?.[0];

    if (!file) {
        statusElement.textContent = "Сначала выберите изображение.";
        resultElement.hidden = true;
        return;
    }

    submitButton.disabled = true;
    submitButton.textContent = "Определяем...";
    statusElement.textContent = "Идёт обработка изображения...";
    resultElement.hidden = true;

    try {
        const formData = new FormData();
        formData.append("file", file);

        const response = await fetch(API_URL, {
            method: "POST",
            body: formData,
        });

        if (!response.ok) {
            throw new Error(`Ошибка сервера: ${response.status}`);
        }

        const data = await response.json();

        breedValue.textContent = data.breed ?? "Неизвестно";
        confidenceValue.textContent = `${((data.confidence ?? 0) * 100).toFixed(1)}%`;
        modeValue.textContent = data.source === "mock" ? "Мок-ответ" : "ML модель";

        resultElement.hidden = false;
        statusElement.textContent = "Результат успешно получен.";
    } catch (error) {
        console.error(error);
        statusElement.textContent = "Не удалось получить результат. Проверь backend и CORS.";
        resultElement.hidden = true;
    } finally {
        submitButton.disabled = false;
        submitButton.textContent = "Определить породу";
    }
});

window.addEventListener("beforeunload", () => {
    clearPreviewUrl();
});