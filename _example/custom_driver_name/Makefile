TARGET = custom_driver_name
ifeq ($(OS),Windows_NT)
TARGET := $(TARGET).exe
endif

all : $(TARGET)

$(TARGET) : main.go
	go build -ldflags="-X 'github.com/thinkgos/go-sqlcipher.driverName=my-sqlite3'"

clean :
	rm -f $(TARGET)
