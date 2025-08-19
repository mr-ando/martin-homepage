import { Provider } from '../components/ui/provider';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import RTCApp from './RTCApp';
const queryClient = new QueryClient();

export default function RTCWrapper() {
  return (
    <div className="flex-1">
      <Provider>
        <QueryClientProvider client={queryClient}>
          <RTCApp />
        </QueryClientProvider>
      </Provider>
    </div>
  );
}
