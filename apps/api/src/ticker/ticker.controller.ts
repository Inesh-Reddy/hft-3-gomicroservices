import { Controller, Get, Query, Inject, Sse } from '@nestjs/common';
import { ClientGrpc } from '@nestjs/microservices';
import { Observable } from 'rxjs';
import {
  TickerServiceClient,
  TickerUpdate,
  TICKER_SERVICE_NAME,
} from '../proto/ticker';

@Controller('ticker')
export class TickerController {
  private tickerService: TickerServiceClient;

  constructor(@Inject('TICKER_PACKAGE') private readonly client: ClientGrpc) {
    this.tickerService =
      this.client.getService<TickerServiceClient>(TICKER_SERVICE_NAME);
  }

  @Get()
  @Sse()
  getTicker(@Query('symbol') symbol: string): Observable<TickerUpdate> {
    console.log('Received symbol:', symbol); // Debug log
    if (!symbol) {
      throw new Error('Symbol is required');
    }
    const data = this.tickerService.streamTicker({ symbol });
    console.log(data);
    console.log(data.pipe());
    return data.pipe();
  }
}
