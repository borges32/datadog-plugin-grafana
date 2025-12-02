import React, { ChangeEvent } from 'react';
import { InlineField, TextArea } from '@grafana/ui';
import { QueryEditorProps } from '@grafana/data';
import { DataSource } from '../datasource';
import { MyDataSourceOptions, MyQuery } from '../types';

type Props = QueryEditorProps<DataSource, MyQuery, MyDataSourceOptions>;

export function QueryEditor({ query, onChange, onRunQuery }: Props) {
  const onQueryChange = (event: ChangeEvent<HTMLTextAreaElement>) => {
    onChange({ ...query, query: event.target.value });
  };

  const { query: queryText } = query;

  return (
    <div className="gf-form">
      <InlineField 
        label="Query" 
        labelWidth={16} 
        tooltip="Datadog query (e.g., avg:processor.time{host:XX-XX-XX-XX} by {host,instance})"
        grow
      >
        <TextArea 
          onChange={onQueryChange} 
          onBlur={onRunQuery}
          value={queryText || ''} 
          rows={3}
          placeholder="avg:processor.time{host:XX-XX-XX-XX} by {host,instance}"
        />
      </InlineField>
    </div>
  );
}
