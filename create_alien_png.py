#!/usr/bin/env python3
"""
Create a simple alien face PNG for the bouncing balls game.
This creates a green alien face with large black eyes and a small mouth.
"""

from PIL import Image, ImageDraw
import sys

def create_alien_face(size=100):
    # Create a new image with transparent background
    img = Image.new('RGBA', (size, size), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)

    # Draw alien head (green oval)
    head_margin = 5
    head_color = (100, 255, 100, 255)  # Bright green
    draw.ellipse([head_margin, head_margin, size-head_margin, size-head_margin],
                 fill=head_color, outline=(50, 200, 50, 255), width=2)

    # Draw large alien eyes (black)
    eye_size = size // 6
    left_eye_x = size // 3 - eye_size // 2
    right_eye_x = 2 * size // 3 - eye_size // 2
    eye_y = size // 3

    # Left eye
    draw.ellipse([left_eye_x, eye_y, left_eye_x + eye_size, eye_y + eye_size],
                 fill=(0, 0, 0, 255))

    # Right eye
    draw.ellipse([right_eye_x, eye_y, right_eye_x + eye_size, eye_y + eye_size],
                 fill=(0, 0, 0, 255))

    # Draw small mouth
    mouth_width = size // 8
    mouth_height = size // 16
    mouth_x = size // 2 - mouth_width // 2
    mouth_y = 2 * size // 3

    draw.ellipse([mouth_x, mouth_y, mouth_x + mouth_width, mouth_y + mouth_height],
                 fill=(50, 50, 50, 255))

    return img

def main():
    try:
        # Create alien face
        alien_img = create_alien_face(100)

        # Save as PNG
        alien_img.save('alien.png', 'PNG')
        print("✅ Created alien.png successfully!")

    except ImportError:
        print("❌ PIL (Pillow) not installed. Please install with: pip install Pillow")
        print("Or you can manually create an alien.png file and place it in the project directory.")
        sys.exit(1)
    except Exception as e:
        print(f"❌ Error creating alien.png: {e}")
        sys.exit(1)

if __name__ == "__main__":
    main()