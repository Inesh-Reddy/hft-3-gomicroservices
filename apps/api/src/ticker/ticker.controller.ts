import { Controller, Get, Query, Inject, OnModuleInit } from '@nestjs/common';
import { ClientGrpc } from '@nestjs/microservices';
import { Observable } from 'rxjs';
import {
  TickerServiceClient,
  TickerUpdate,
  TICKER_SERVICE_NAME,
} from '../proto/ticker'; // Matches src/proto

@Controller('ticker')
export class TickerController implements OnModuleInit {
  private tickerService: TickerServiceClient;

  constructor(@Inject('TICKER_PACKAGE') private readonly client: ClientGrpc) {}

  onModuleInit() {
    this.tickerService =
      this.client.getService<TickerServiceClient>(TICKER_SERVICE_NAME);
  }

  @Get()
  getTicker(@Query('symbol') symbol: string): Observable<TickerUpdate> {
    return this.tickerService.streamTicker({ symbol });
  }
}
