BUILD_PATH ?= .
ETC_DIR = ./etc
SCRIPTS_DIR = ./scripts
APP_UI_DIR = ./app-ui
SERVICE_NAME = topology
SERVICE_BIN = $(SERVICE_NAME)
SERVICE_CONFIG_FILE = $(ETC_DIR)/$(SERVICE_NAME).yaml
SERVICE_TGZ = $(SERVICE_NAME).tgz
TEMP_INSTALL_DIR = ./.install
INSTALL_DIR = ./install

GLIMPSE_DIR = ../../SpirentOrion/glimpse
GLIMPSE_BIN = $(GLIMPSE_DIR)/bin/glimpse

.PHONY: all install clean

all: $(SERVICE_BIN)

vet:
	go tool vet -all -composites=false -shadow=true .

$(SERVICE_BIN): $(BIN_DIR)/
	#cd $(GLIMPSE_DIR) && $(MAKE)
	go build -o $(SERVICE_BIN) ./.

$(INSTALL_DIR)/$(SERVICE_TGZ): $(SERVICE_BIN) $(SERVICE_CONFIG_FILE)
	mkdir -p $(TEMP_INSTALL_DIR) && \
	mkdir -p $(TEMP_INSTALL_DIR)/bin && \
	mkdir -p $(TEMP_INSTALL_DIR)/logs && \
	cp -R $(SCRIPTS_DIR)/* $(TEMP_INSTALL_DIR)/bin && \
	cp $(SERVICE_BIN) $(TEMP_INSTALL_DIR)/bin && \
	cp $(GLIMPSE_BIN) $(TEMP_INSTALL_DIR)/bin && \
	cp -R $(APP_UI_DIR) $(TEMP_INSTALL_DIR) && \
	cp -R $(ETC_DIR) $(TEMP_INSTALL_DIR) && \
	cd $(TEMP_INSTALL_DIR) && \
	tar -czf $(SERVICE_TGZ) * && \
	cd .. && \
	mkdir -p $(INSTALL_DIR) && \
	mv $(TEMP_INSTALL_DIR)/$(SERVICE_TGZ) $(INSTALL_DIR) && \
	rm -rf $(TEMP_INSTALL_DIR)

install: $(INSTALL_DIR)/$(SERVICE_TGZ)

clean:
	cd $(BUILD_PATH) && \
	rm -rf $(TEMP_INSTALL_DIR) && \
	rm -rf $(INSTALL_DIR) && \
	go clean -i ./.
