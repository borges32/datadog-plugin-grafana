import React, { ChangeEvent } from 'react';
import { InlineField, Input, SecretInput } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { MyDataSourceOptions, MySecureJsonData } from '../types';

interface Props extends DataSourcePluginOptionsEditorProps<MyDataSourceOptions> {}

export function ConfigEditor(props: Props) {
  const { onOptionsChange, options } = props;
  
  const onURLChange = (event: ChangeEvent<HTMLInputElement>) => {
    const jsonData = {
      ...options.jsonData,
      url: event.target.value,
    };
    onOptionsChange({ ...options, jsonData });
  };

  // Secure field (only sent to the backend)
  const onAPIKeyChange = (event: ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      secureJsonData: {
        ...options.secureJsonData,
        apiKey: event.target.value,
      },
    });
  };

  const onResetAPIKey = () => {
    onOptionsChange({
      ...options,
      secureJsonFields: {
        ...options.secureJsonFields,
        apiKey: false,
      },
      secureJsonData: {
        ...options.secureJsonData,
        apiKey: '',
      },
    });
  };

  const onApplicationKeyChange = (event: ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      secureJsonData: {
        ...options.secureJsonData,
        applicationKey: event.target.value,
      },
    });
  };

  const onResetApplicationKey = () => {
    onOptionsChange({
      ...options,
      secureJsonFields: {
        ...options.secureJsonFields,
        applicationKey: false,
      },
      secureJsonData: {
        ...options.secureJsonData,
        applicationKey: '',
      },
    });
  };

  const { jsonData, secureJsonFields } = options;
  const secureJsonData = (options.secureJsonData || {}) as MySecureJsonData;

  return (
    <div className="gf-form-group">
      <InlineField label="URL" labelWidth={12} tooltip="Datadog API URL (e.g., https://api.us3.datadoghq.com/)">
        <Input
          onChange={onURLChange}
          value={jsonData.url || ''}
          placeholder="https://api.us3.datadoghq.com/"
          width={40}
        />
      </InlineField>
      <InlineField label="DD-API-KEY" labelWidth={12} tooltip="Datadog API Key">
        <SecretInput
          isConfigured={(secureJsonFields && secureJsonFields.apiKey) as boolean}
          value={secureJsonData.apiKey || ''}
          placeholder="Datadog API Key"
          width={40}
          onReset={onResetAPIKey}
          onChange={onAPIKeyChange}
        />
      </InlineField>
      <InlineField label="DD-APPLICATION-KEY" labelWidth={12} tooltip="Datadog Application Key">
        <SecretInput
          isConfigured={(secureJsonFields && secureJsonFields.applicationKey) as boolean}
          value={secureJsonData.applicationKey || ''}
          placeholder="Datadog Application Key"
          width={40}
          onReset={onResetApplicationKey}
          onChange={onApplicationKeyChange}
        />
      </InlineField>
    </div>
  );
}
