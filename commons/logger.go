package commons

import (
	"context"

	"errors"
	"reflect"
	"sync"
	"time"

	newrelic "github.com/newrelic/go-agent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var loggerOnce sync.Once
var goflowLogger *Logger

func EnableLogging() *Logger {
	loggerOnce.Do(func() {
		goflowLogger = &Logger{
			additionalFields: make(map[string]interface{}),
		}
	})

	return goflowLogger
}

func GetLogger() *Logger {
	return goflowLogger
}

type LoggingField struct {
	key         *string
	fieldType   LoggingFieldType
	boolData    *bool
	integerData *int
	stringData  *string
	errorData   error
	anyData     interface{}
}

type logMessage struct {
	ctx           context.Context
	message       *string
	logLevel      LogLevel
	loggingFields []*LoggingField
}

type asyncLogConfig struct {
	channelCapacity int
	waitGrp         *sync.WaitGroup
	logChannel      chan *logMessage
}

func (asyncLogConfig *asyncLogConfig) MyType() string {
	return reflect.TypeOf(asyncLogConfig).String()
}

func (asyncLogConfig *asyncLogConfig) MyName() string {
	return ""
}

func (asyncLogConfig *asyncLogConfig) Describe() map[string]interface{} {
	descriptionMap := make(map[string]interface{})
	descriptionMap["my type"] = asyncLogConfig.MyType()
	descriptionMap["my name"] = asyncLogConfig.MyName()
	descriptionMap["log channel capacity"] = asyncLogConfig.channelCapacity
	return descriptionMap
}

type Logger struct {
	async            bool
	timeKeyName      string
	levelKeyName     string
	messageKeyName   string
	zapConfig        *zap.Config
	coreLogger       *zap.Logger
	asyncLogConfig   *asyncLogConfig
	fileDescriptors  []FileDescriptor
	environment      GoflowRuntimeEnvironment
	additionalFields map[string]interface{}
}

func (logger *Logger) MyType() string {
	return reflect.TypeOf(logger).String()
}

func (logger *Logger) MyName() string {
	return "GoFlowLogger"
}

func (logger *Logger) Describe() map[string]interface{} {
	descriptionMap := make(map[string]interface{})
	descriptionMap["my type"] = logger.MyType()
	descriptionMap["my name"] = logger.MyName()
	descriptionMap["async enabled"] = logger.async
	if logger.async {
		descriptionMap["async config"] = logger.asyncLogConfig.Describe()
	}
	descriptionMap["core logger type"] = reflect.TypeOf(logger.coreLogger).String()
	descriptionMap["output file descriptor"] = logger.fileDescriptors
	descriptionMap["logger environment"] = logger.environment.String()
	return descriptionMap
}

func (logger *Logger) IsCoreLoggerSet() bool {
	return logger.coreLogger != nil
}

func (logger *Logger) setCoreLogger(coreLogger *zap.Logger) {
	logger.coreLogger = coreLogger
}

func (logger *Logger) AddFileDescriptor(fd FileDescriptor) {
	logger.fileDescriptors = append(logger.fileDescriptors, fd)
}

func (logger *Logger) SetLogEnvironment(env GoflowRuntimeEnvironment) {
	logger.environment = env
}

func (logger *Logger) SetZapConfig(cfg *zap.Config) {
	logger.zapConfig = cfg
}

func (logger *Logger) AddAdditionalFields(key string, value interface{}) {
	logger.additionalFields[key] = value
}

func (logger *Logger) SetTimeKey(keyName string) {
	logger.timeKeyName = keyName
}

func (logger *Logger) SetLevelKey(keyName string) {
	logger.levelKeyName = keyName
}

func (logger *Logger) SetMsgKey(keyName string) {
	logger.messageKeyName = keyName
}

func (logger *Logger) DoAsyncLogging(capacity int, wg *sync.WaitGroup) error {
	if wg == nil {
		return errors.New("nil wait group")
	}

	logger.async = true
	logger.asyncLogConfig = &asyncLogConfig{
		waitGrp:         wg,
		channelCapacity: capacity,
	}

	return nil
}

func (logger *Logger) Build() error {
	if (len(logger.fileDescriptors) == 0) ||
		(logger.environment == UndefinedRuntimeEnvironment) {
		return errors.New("logger configuration not set properly")
	}

	if logger.zapConfig == nil {
		var loggerConfig zap.Config
		switch logger.environment {
		case Development:
			loggerConfig = zap.NewDevelopmentConfig()
		case PreProduction:
			loggerConfig = zap.NewProductionConfig()
		case ProductionPP:
			loggerConfig = zap.NewProductionConfig()
		case Production:
			loggerConfig = zap.NewProductionConfig()
		default:
			return errors.New("unknown logging environment set")
		}

		if logger.levelKeyName == "" {
			logger.levelKeyName = DefaultLoggingLevelKeyName
		}

		if logger.timeKeyName == "" {
			logger.timeKeyName = DefaultLoggingTimeKeyName
		}

		if logger.messageKeyName == "" {
			logger.messageKeyName = DefaultLoggingMessageKeyName
		}

		loggerConfig.EncoderConfig.LevelKey = logger.levelKeyName
		loggerConfig.EncoderConfig.TimeKey = logger.timeKeyName
		loggerConfig.EncoderConfig.MessageKey = logger.messageKeyName
		loggerConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		loggerConfig.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

		// suppressing these
		loggerConfig.EncoderConfig.CallerKey = ""
		loggerConfig.EncoderConfig.EncodeCaller = nil
		loggerConfig.EncoderConfig.StacktraceKey = ""

		loggerConfig.InitialFields = logger.additionalFields
		for _, fileDescriptor := range logger.fileDescriptors {
			loggerConfig.OutputPaths = []string{fileDescriptor.String()}
		}

		logger.zapConfig = &loggerConfig
	}

	zapLogger, zapLoggerBuildErr := logger.zapConfig.Build()
	if zapLoggerBuildErr != nil {
		return WrapError(zapLoggerBuildErr, "fatal error building zap Logger")
	}

	if logger.async {
		if logger.asyncLogConfig.waitGrp == nil {
			return errors.New("waitgroup nil for async logging")
		}

		if logger.asyncLogConfig.channelCapacity <= 0 {
			logger.asyncLogConfig.channelCapacity = DefaultAsyncLoggerChannelCapacity
		}

		logger.asyncLogConfig.logChannel = make(chan *logMessage, logger.asyncLogConfig.channelCapacity)
		go logger.asyncLogging()
		logger.asyncLogConfig.waitGrp.Add(1)
	}

	logger.setCoreLogger(zapLogger)

	return nil
}

func (logger *Logger) Bool(key string, val bool) *LoggingField {
	return &LoggingField{key: &key, fieldType: Bool, boolData: &val}
}

func (logger *Logger) Int(key string, val int) *LoggingField {
	return &LoggingField{key: &key, fieldType: Int, integerData: &val}
}

func (logger *Logger) Error(val error) *LoggingField {
	return &LoggingField{key: nil, fieldType: Error, errorData: val}
}

func (logger *Logger) String(key string, val string) *LoggingField {
	return &LoggingField{key: &key, fieldType: String, stringData: &val}
}

func (logger *Logger) Interface(key string, val interface{}) *LoggingField {
	return &LoggingField{key: &key, fieldType: Interface, anyData: &val}
}

// do not use this in production
func (logger *Logger) Description(enitity Describable) *LoggingField {
	name := enitity.MyName()
	desc := enitity.Describe()
	return &LoggingField{key: &name, fieldType: Interface, anyData: &desc}
}

func (logger *Logger) LogInfo(ctx context.Context, message string, fields ...*LoggingField) {
	if logger.async {
		message := logMessage{
			ctx:           ctx,
			message:       &message,
			logLevel:      INFO,
			loggingFields: fields,
		}
		logger.asyncLogConfig.logChannel <- &message
	} else {
		logger.syncLogging(ctx, message, INFO, fields...)
	}

	return
}

func (logger *Logger) LogDebug(ctx context.Context, message string, fields ...*LoggingField) {
	// if logger.environment > enigma_constants.DEVLOPMENT {
	// 	return
	// }

	if logger.async {
		message := logMessage{
			ctx:           ctx,
			message:       &message,
			logLevel:      DEBUG,
			loggingFields: fields,
		}
		logger.asyncLogConfig.logChannel <- &message
	} else {
		logger.syncLogging(ctx, message, DEBUG, fields...)
	}

	return
}

func (logger *Logger) LogWarn(ctx context.Context, message string, fields ...*LoggingField) {
	if logger.async {
		message := logMessage{
			ctx:           ctx,
			message:       &message,
			logLevel:      WARN,
			loggingFields: fields,
		}
		logger.asyncLogConfig.logChannel <- &message
	} else {
		logger.syncLogging(ctx, message, WARN, fields...)
	}

	return
}

func (logger *Logger) LogError(ctx context.Context, err error, fields ...*LoggingField) {
	var errStr string
	if err != nil {
		errStr = err.Error()
	}

	if logger.async {
		message := logMessage{
			ctx:           ctx,
			message:       &errStr,
			logLevel:      ERROR,
			loggingFields: fields,
		}
		logger.asyncLogConfig.logChannel <- &message
	} else {
		logger.syncLogging(ctx, errStr, ERROR, fields...)
	}

	return
}

func (logger *Logger) LogDPanic(ctx context.Context, message string, fields ...*LoggingField) {
	if logger.async {
		message := logMessage{
			ctx:           ctx,
			message:       &message,
			logLevel:      DPANIC,
			loggingFields: fields,
		}
		logger.asyncLogConfig.logChannel <- &message
	} else {
		logger.syncLogging(ctx, message, DPANIC, fields...)
	}

	return
}

func (logger *Logger) LogPanic(ctx context.Context, message string, fields ...*LoggingField) {
	if logger.async {
		message := logMessage{
			ctx:           ctx,
			message:       &message,
			logLevel:      DPANIC,
			loggingFields: fields,
		}
		logger.asyncLogConfig.logChannel <- &message
	} else {
		logger.syncLogging(ctx, message, DPANIC, fields...)
	}

	return
}

func (logger *Logger) LogFatal(ctx context.Context, message string, fields ...*LoggingField) {
	if logger.async {
		message := logMessage{
			ctx:           ctx,
			message:       &message,
			logLevel:      DPANIC,
			loggingFields: fields,
		}
		logger.asyncLogConfig.logChannel <- &message
	} else {
		logger.syncLogging(ctx, message, DPANIC, fields...)
	}

	return
}

func (logger *Logger) syncLogging(ctx context.Context, message string, logLevel LogLevel, fields ...*LoggingField) {
	uuid := ""
	if ctx != nil {
		value := ctx.Value(*UUIDGoFlowContextKey())
		uuid, _ = value.(string)
	}

	zapFieldsArr := make([]zap.Field, len(fields))
	for index, field := range fields {
		switch field.fieldType {
		case Bool:
			zapFieldsArr[index] = zap.Bool(*field.key, *field.boolData)
		case Int:
			zapFieldsArr[index] = zap.Int(*field.key, *field.integerData)
		case String:
			zapFieldsArr[index] = zap.String(*field.key, *field.stringData)
		case Error:
			zapFieldsArr[index] = zap.Error(field.errorData)
		case Interface:
			zapFieldsArr[index] = zap.Reflect(*field.key, field.anyData)
		default:
		}
	}
	zapFieldsArr = append(zapFieldsArr, zap.String(DefaultLoggingUUIDKeyName, uuid))

	switch logLevel {
	case INFO:
		logger.coreLogger.Info(message, zapFieldsArr...)
	case DEBUG:
		logger.coreLogger.Debug(message, zapFieldsArr...)
	case WARN:
		logger.coreLogger.Warn(message, zapFieldsArr...)
	case ERROR:
		logger.coreLogger.Error(message, zapFieldsArr...)
	case DPANIC:
		logger.coreLogger.DPanic(message, zapFieldsArr...)
	case PANIC:
		logger.coreLogger.Panic(message, zapFieldsArr...)
	case FATAL:
		logger.coreLogger.Fatal(message, zapFieldsArr...)
	}

	return
}

func (logger *Logger) asyncLogging() {
	defer logger.asyncLogConfig.waitGrp.Done()

	var ticker *time.Ticker
	var newRelicApp newrelic.Application
	var logTransactionName string
	var logTransactionLogTime int
	var logConsumerRoutine newrelic.Transaction

	newRMonitor := GetNewRelicMonitoring()
	if newRMonitor != nil {
		newRelicApp = newRMonitor.GetNewRelicApplication()
		logTransactionName = newRMonitor.GetLoggerTransactionName()
		logTransactionLogTime = newRMonitor.GetLoggerTransactionLogTime()

		logConsumerRoutine = newRelicApp.StartTransaction(logTransactionName, nil, nil)
	} else {
		// give a positive value initially, ticker will stop automatically
		logTransactionLogTime = 1
	}

	ticker = time.NewTicker(time.Duration(logTransactionLogTime) * time.Minute)
	for {
		select {
		case logMsg := <-logger.asyncLogConfig.logChannel:
			if logMsg == nil {
				continue
			}

			uuid := ""
			if logMsg.ctx != nil {
				uuid = logMsg.ctx.Value(*UUIDGoFlowContextKey()).(string)
			}

			zapFieldsArr := make([]zap.Field, len(logMsg.loggingFields))
			for index, field := range logMsg.loggingFields {
				switch field.fieldType {
				case Bool:
					zapFieldsArr[index] = zap.Bool(*field.key, *field.boolData)
				case Int:
					zapFieldsArr[index] = zap.Int(*field.key, *field.integerData)
				case String:
					zapFieldsArr[index] = zap.String(*field.key, *field.stringData)
				case Error:
					zapFieldsArr[index] = zap.Error(field.errorData)
				case Interface:
					zapFieldsArr[index] = zap.Reflect(*field.key, field.anyData)
				default:
				}
			}
			zapFieldsArr = append(zapFieldsArr, zap.String(DefaultLoggingUUIDKeyName, uuid))
			switch logMsg.logLevel {
			case INFO:
				logger.coreLogger.Info(*logMsg.message, zapFieldsArr...)
			case DEBUG:
				logger.coreLogger.Debug(*logMsg.message, zapFieldsArr...)
			case WARN:
				logger.coreLogger.Warn(*logMsg.message, zapFieldsArr...)
			case ERROR:
				logger.coreLogger.Error(*logMsg.message, zapFieldsArr...)
			case DPANIC:
				logger.coreLogger.DPanic(*logMsg.message, zapFieldsArr...)
			case PANIC:
				logger.coreLogger.Panic(*logMsg.message, zapFieldsArr...)
			case FATAL:
				logger.coreLogger.Fatal(*logMsg.message, zapFieldsArr...)
			}

		case <-ticker.C:
			if logConsumerRoutine != nil {
				logConsumerRoutine.End()
				logConsumerRoutine = newRelicApp.StartTransaction(logTransactionName, nil, nil)
			} else {
				if ticker != nil {
					ticker.Stop()
				}
			}
		}
	}
}
