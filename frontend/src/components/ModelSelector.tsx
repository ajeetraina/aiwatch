import { useState, useEffect } from 'react';
import { DockerModel } from '../types';

interface ModelSelectorProps {
  selectedModel: string;
  onSelectModel: (model: string) => void;
}

export const ModelSelector = ({ selectedModel, onSelectModel }: ModelSelectorProps) => {
  const [models, setModels] = useState<DockerModel[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isOpen, setIsOpen] = useState(false);

  // Fetch available models on component mount
  useEffect(() => {
    const fetchModels = async () => {
      try {
        setIsLoading(true);
        setError(null);
        
        const response = await fetch('http://localhost:8080/models');
        if (!response.ok) {
          throw new Error(`Failed to fetch models: ${response.statusText}`);
        }
        
        const data = await response.json();
        setModels(data);
      } catch (err) {
        console.error('Error fetching models:', err);
        setError(err instanceof Error ? err.message : 'Failed to fetch models');
      } finally {
        setIsLoading(false);
      }
    };

    fetchModels();
  }, []);

  // Get the selected model details
  const getSelectedModelDetails = (): DockerModel | null => {
    const found = models.find(model => model.name === selectedModel);
    return found || null;
  };

  // Handle model selection
  const handleSelectModel = (modelName: string) => {
    onSelectModel(modelName);
    setIsOpen(false);
  };

  // Get displaying name for the model (short version)
  const getDisplayName = (model: DockerModel): string => {
    // Extract the model name without repository prefix for cleaner display
    const nameParts = model.name.split('/');
    return nameParts[nameParts.length - 1];
  };

  return (
    <div className="relative">
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="flex items-center justify-between w-56 px-3 py-2 text-sm font-medium bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm hover:bg-gray-50 dark:hover:bg-gray-700 focus:outline-none"
        aria-haspopup="listbox"
        aria-expanded={isOpen}
        data-testid="model-selector-button"
      >
        <span className="flex items-center">
          {isLoading ? (
            <span className="text-gray-500">Loading models...</span>
          ) : error ? (
            <span className="text-red-500">Error loading models</span>
          ) : getSelectedModelDetails() ? (
            <>
              <span className="ml-2 truncate">{getDisplayName(getSelectedModelDetails()!)}</span>
              <span className="ml-1 text-xs text-gray-500 truncate">
                ({getSelectedModelDetails()?.parameters})
              </span>
            </>
          ) : (
            <span className="text-gray-700 dark:text-gray-300">Select a model</span>
          )}
        </span>
        <svg
          className="w-5 h-5 ml-2 -mr-1 text-gray-400"
          viewBox="0 0 20 20"
          fill="currentColor"
          aria-hidden="true"
        >
          <path
            fillRule="evenodd"
            d="M5.23 7.21a.75.75 0 011.06.02L10 11.168l3.71-3.938a.75.75 0 111.08 1.04l-4.25 4.5a.75.75 0 01-1.08 0l-4.25-4.5a.75.75 0 01.02-1.06z"
            clipRule="evenodd"
          />
        </svg>
      </button>

      {isOpen && (
        <div className="absolute mt-1 w-64 max-h-56 overflow-auto rounded-md bg-white dark:bg-gray-800 shadow-lg border border-gray-300 dark:border-gray-700 z-10">
          <ul
            className="py-1 text-sm text-gray-700 dark:text-gray-200"
            role="listbox"
            aria-labelledby="model-selector-button"
          >
            {isLoading ? (
              <li className="px-4 py-2 text-gray-500">Loading models...</li>
            ) : error ? (
              <li className="px-4 py-2 text-red-500">{error}</li>
            ) : models.length === 0 ? (
              <li className="px-4 py-2 text-gray-500">No models available</li>
            ) : (
              models.map((model) => (
                <li
                  key={model.modelId}
                  className={`px-4 py-2 hover:bg-gray-100 dark:hover:bg-gray-700 cursor-pointer ${
                    model.name === selectedModel
                      ? 'bg-gray-100 dark:bg-gray-700 font-medium'
                      : ''
                  }`}
                  role="option"
                  aria-selected={model.name === selectedModel}
                  onClick={() => handleSelectModel(model.name)}
                >
                  <div className="flex flex-col">
                    <span className="font-medium">{getDisplayName(model)}</span>
                    <span className="text-xs text-gray-500 dark:text-gray-400">
                      {model.parameters} • {model.quantization} • {model.architecture}
                    </span>
                    <span className="text-xs text-gray-500 dark:text-gray-400">
                      Size: {model.size} • Created: {model.created}
                    </span>
                  </div>
                </li>
              ))
            )}
          </ul>
        </div>
      )}
    </div>
  );
};

export default ModelSelector;
