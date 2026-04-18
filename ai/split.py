import os
import shutil
import random
from pathlib import Path

# Фиксируем random seed для воспроизводимости
random.seed(42)

# Пути
BASE_DIR = Path("./data/catbreeds")
TRAIN_DIR = Path("./data/train")
VAL_DIR = Path("./data/val")

# Соотношение train/val
TRAIN_RATIO = 0.8

# Минимальное количество фото для валидации
MIN_VAL_IMAGES = 5


def split_and_copy(breed: str, breed_path: Path):
    """
    Случайно разделяет изображения одной породы на train и val.
    """
    # Все изображения породы (поддерживаемые форматы)
    image_extensions = {'.jpg', '.jpeg', '.png', '.gif', '.bmp', '.webp'}
    images = [f for f in breed_path.iterdir() 
              if f.is_file() and f.suffix.lower() in image_extensions]
    
    if not images:
        print(f"  ⚠️  {breed}: нет изображений")
        return None, None
    
    # Перемешиваем все изображения
    random.shuffle(images)
    
    # Вычисляем точку разделения
    total_images = len(images)
    split_idx = int(total_images * TRAIN_RATIO)
    
    # Корректируем для малого количества
    if total_images <= MIN_VAL_IMAGES:
        # Если фото мало, кладём всё в train
        train_images = images
        val_images = []
        print(f"  ⚠️  {breed}: всего {total_images} фото, всё в train")
    else:
        # Гарантируем хотя бы 1 фото в val если возможно
        if split_idx == total_images:
            split_idx = total_images - 1
        
        train_images = images[:split_idx]
        val_images = images[split_idx:]
    
    # Создаём папки
    train_breed_dir = TRAIN_DIR / breed
    val_breed_dir = VAL_DIR / breed
    train_breed_dir.mkdir(parents=True, exist_ok=True)
    val_breed_dir.mkdir(parents=True, exist_ok=True)
    
    # Копируем файлы
    for img_path in train_images:
        shutil.copy2(img_path, train_breed_dir / img_path.name)
        
    for img_path in val_images:
        shutil.copy2(img_path, val_breed_dir / img_path.name)
    
    return len(train_images), len(val_images)


def main():
    print("=" * 60)
    print("РАЗДЕЛЕНИЕ ДАТАСЕТА НА TRAIN / VAL")
    print("=" * 60)
    print(f"Исходная папка: {BASE_DIR}")
    print(f"Train: {TRAIN_DIR}")
    print(f"Val:   {VAL_DIR}")
    print(f"Соотношение: {TRAIN_RATIO:.0%} / {1-TRAIN_RATIO:.0%}")
    print("=" * 60)
    
    # Очищаем train и val если они существуют
    if TRAIN_DIR.exists():
        print(f"🗑️  Удаляю старый train...")
        shutil.rmtree(TRAIN_DIR)
    if VAL_DIR.exists():
        print(f"🗑️  Удаляю старый val...")
        shutil.rmtree(VAL_DIR)
    
    TRAIN_DIR.mkdir(parents=True, exist_ok=True)
    VAL_DIR.mkdir(parents=True, exist_ok=True)
    
    # Получаем список всех пород
    breeds = sorted([d.name for d in BASE_DIR.iterdir() if d.is_dir()])
    
    print(f"\n📁 Найдено пород: {len(breeds)}")
    print("-" * 60)
    
    total_train = 0
    total_val = 0
    
    for breed in breeds:
        breed_path = BASE_DIR / breed
        train_count, val_count = split_and_copy(breed, breed_path)
        
        if train_count is not None:
            total_train += train_count
            total_val += val_count
            val_str = f"{val_count:>4}" if val_count > 0 else "   0"
            print(f"  ✓ {breed:<25} train: {train_count:>4} | val: {val_str}")
    
    print("-" * 60)
    print(f"ВСЕГО:")
    print(f"  Train: {total_train} изображений")
    print(f"  Val:   {total_val} изображений")
    print(f"  Total: {total_train + total_val} изображений")
    print("=" * 60)
    print("✅ Готово!")


if __name__ == "__main__":
    main()
