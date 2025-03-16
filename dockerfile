# Use an official Python image
FROM python:3.10

# Set environment variables
ENV PYTHONUNBUFFERED=1
ENV PYTHONDONTWRITEBYTECODE=1

# Set working directory
WORKDIR /app

# Copy application files
COPY . /app

# Copy pre-downloaded dependencies
COPY wheelhouse /app/wheelhouse

# Install dependencies offline
RUN pip install --no-index --find-links=/app/wheelhouse -r requirements.txt

# Expose port (optional)
EXPOSE 8000

# Start Django application
CMD ["gunicorn", "--bind", "0.0.0.0:8000", "your_project.wsgi:application"]
