declare module '@wangeditor/editor-for-vue' {
  import type { DefineComponent } from 'vue';
  import type { IDomEditor } from '@wangeditor/editor';

  export const Editor: DefineComponent<{
    modelValue?: string;
    defaultConfig?: any;
    mode?: string;
    onCreated?: (editor: IDomEditor) => void;
  }>;

  export const Toolbar: DefineComponent<{
    editor?: IDomEditor;
    defaultConfig?: any;
    mode?: string;
  }>;
}
