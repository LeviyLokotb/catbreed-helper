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
        statusElement.textContent = "Пожалуйста, выберите именно изображение.";
        resetPreview();
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
    statusElement.textContent = "Идёт обработка изображения...";
    resultElement.hidden = true;

    try {
        const formData = new FormData();
        formData.append("file", file);

        const response = await fetch("http://127.0.0.1:8000/predict", {
            method: "POST",
            body: formData,
        });

        if (!response.ok) {
            throw new Error("Ошибка сервера");
        }

        const data = await response.json();

        breedValue.textContent = data.breed ?? "Неизвестно";
        confidenceValue.textContent = `${((data.confidence ?? 0) * 100).toFixed(1)}%`;
        modeValue.textContent = data.mode === "mock" ? "Мок-ответ" : "ML модель";

        resultElement.hidden = false;
        statusElement.textContent = "Результат успешно получен.";
    } catch (error) {
        console.error(error);
        statusElement.textContent = "Не удалось получить результат. Проверь backend.";
        resultElement.hidden = true;
    } finally {
        submitButton.disabled = false;
    }
});