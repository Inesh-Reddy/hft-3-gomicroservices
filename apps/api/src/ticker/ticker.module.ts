import { Module } from '@nestjs/common';
import { TickerController } from './ticker.controller';

@Module({
  controllers: [TickerController],
})
export class TickerModule {}
